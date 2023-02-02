//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package replica

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate/entities/additional"
	"github.com/weaviate/weaviate/entities/search"
	"github.com/weaviate/weaviate/entities/storobj"
)

// ErrConsistencyLevel consistency level cannot be achieved
var ErrConsistencyLevel = errors.New("cannot achieve consistency level")

type (
	// senderReply represent the data received from a sender
	senderReply[T any] struct {
		sender     string // hostname of the sender
		Version    int64  // sender's current version of the object
		data       T      // the data sent by the sender
		UpdateTime int64  // sender's current update time
		DigestRead bool
	}
	findOneReply    senderReply[*storobj.Object]
	existReply      senderReply[bool]
	getObjectsReply senderReply[[]*storobj.Object]
)

// Finder finds replicated objects
type Finder struct {
	RClient            // needed to commit and abort operation
	resolver *resolver // host names of replicas
	class    string
}

// NewFinder constructs a new finder instance
func NewFinder(className string,
	stateGetter shardingState, nodeResolver nodeResolver,
	client RClient,
) *Finder {
	return &Finder{
		class: className,
		resolver: &resolver{
			schema:       stateGetter,
			nodeResolver: nodeResolver,
			class:        className,
		},
		RClient: client,
	}
}

// GetOne gets object which satisfies the giving consistency
func (f *Finder) GetOne(ctx context.Context, l ConsistencyLevel, shard string,
	id strfmt.UUID, props search.SelectProperties, additional additional.Properties,
) (*storobj.Object, error) {
	c := newReadCoordinator[findOneReply](f, shard)
	op := func(ctx context.Context, host string) (findOneReply, error) {
		obj, err := f.FindObject(ctx, host, f.class, shard, id, props, additional)
		return findOneReply{host, -1, obj, 0, false}, err
	}
	replyCh, state, err := c.Fetch(ctx, l, op)
	if err != nil {
		return nil, err
	}
	result := <-readOne(replyCh, state)
	return result.data, result.err
}

// GetOne gets object which satisfies the giving consistency
func (f *Finder) GetOne2(ctx context.Context, l ConsistencyLevel, shard string,
	id strfmt.UUID, props search.SelectProperties, additional additional.Properties,
) (*storobj.Object, error) {
	c := newReadCoordinator[findOneReply](f, shard)
	op := func(ctx context.Context, host string, fullRead bool) (findOneReply, error) {
		obj, err := f.FindObject(ctx, host, f.class, shard, id, props, additional)
		var uTime int64
		if obj != nil {
			uTime = obj.LastUpdateTimeUnix()
		}
		return findOneReply{host, -1, obj, uTime, false}, err
	}
	replyCh, state, err := c.Fetch2(ctx, l, op)
	if err != nil {
		return nil, err
	}
	result := <-f.readOne(replyCh, state)
	return result.data, result.err
}

// Exists checks if an object exists which satisfies the giving consistency
func (f *Finder) Exists(ctx context.Context, l ConsistencyLevel, shard string, id strfmt.UUID) (bool, error) {
	c := newReadCoordinator[existReply](f, shard)
	op := func(ctx context.Context, host string) (existReply, error) {
		obj, err := f.RClient.Exists(ctx, host, f.class, shard, id)
		return existReply{host, -1, obj, 0, false}, err
	}
	replyCh, state, err := c.Fetch(ctx, l, op)
	if err != nil {
		return false, err
	}
	return readOneExists(replyCh, state)
}

// GetAll gets all objects which satisfy the giving consistency
func (f *Finder) GetAll(ctx context.Context, l ConsistencyLevel, shard string,
	ids []strfmt.UUID,
) ([]*storobj.Object, error) {
	c := newReadCoordinator[getObjectsReply](f, shard)
	op := func(ctx context.Context, host string) (getObjectsReply, error) {
		objs, err := f.RClient.MultiGetObjects(ctx, host, f.class, shard, ids)
		return getObjectsReply{host, -1, objs, 0, false}, err
	}
	replyCh, state, err := c.Fetch(ctx, l, op)
	if err != nil {
		return nil, err
	}
	return readAll(replyCh, len(ids), state)
}

// NodeObject gets object from a specific node.
// it is used mainly for debugging purposes
func (f *Finder) NodeObject(ctx context.Context, nodeName, shard string,
	id strfmt.UUID, props search.SelectProperties, additional additional.Properties,
) (*storobj.Object, error) {
	host, ok := f.resolver.NodeHostname(nodeName)
	if !ok || host == "" {
		return nil, fmt.Errorf("cannot resolve node name: %s", nodeName)
	}
	return f.RClient.FindObject(ctx, host, f.class, shard, id, props, additional)
}

func (f *Finder) readOne(ch <-chan simpleResult[findOneReply], st rState) <-chan result[*storobj.Object] {
	// counters tracks the number of votes for each participant
	resultCh := make(chan result[*storobj.Object], 1)
	go func() {
		defer close(resultCh)
		var (
			counters   = make([]objTuple, 0, len(st.Hosts))
			winner     = 0
			max        = 0
			resultSent = false
			contentIdx = 0
		)

		for r := range ch { // len(ch) == st.Level
			resp := r.Response
			if r.Err != nil { // a least one node is not responding
				resultCh <- result[*storobj.Object]{nil, r.Err}
				continue
			}
			if !resp.DigestRead {
				contentIdx = len(counters)
			}
			counters = append(counters, objTuple{resp.sender, resp.UpdateTime, resp.data, 0, nil})

			max = 0
			for i := range counters {
				if counters[i].UTime == resp.UpdateTime {
					counters[i].ack++
				}
				if max < counters[i].ack {
					max = counters[i].ack
					winner = i
				}
				if !resultSent && max >= st.Level {
					resultSent = true
					resultCh <- result[*storobj.Object]{counters[contentIdx].o, nil}
				}
			}
		}
		if len(counters) > 0 && counters[winner].UTime == counters[contentIdx].UTime {
			counters[winner].o = counters[contentIdx].o
		}
		if obj, err := f.repairOne(counters, st, winner); err == nil {
			if !resultSent {
				resultCh <- result[*storobj.Object]{obj, nil}
			}
			return
		}
		if !resultSent {
			var sb strings.Builder
			for i, c := range counters {
				if i != 0 {
					sb.WriteString(", ")
				}
				if c.err != nil {
					fmt.Fprintf(&sb, "%s: %s", c.sender, c.err.Error())
				} else if c.o == nil {
					fmt.Fprintf(&sb, "%s: 0", c.sender)
				} else {
					fmt.Fprintf(&sb, "%s: %d", c.sender, c.o.LastUpdateTimeUnix())
				}
			}
			resultCh <- result[*storobj.Object]{nil, fmt.Errorf("%w %q %s", ErrConsistencyLevel, st.CLevel, sb.String())}
		}
	}()
	return resultCh
}

func (f *Finder) repairOne(counters []objTuple, st rState, winnerIdx int) (*storobj.Object, error) {
	return nil, fmt.Errorf("no majority found")

	for _, x := range counters {
		if x.err != nil {
			return nil, x.err
		}
	}
	// TODO: if winner object is nil we need to tell the node to delete the object
	// The adapter/repos/db/DB.overwriteObjects nil to be adjust to account for nil objects
	vots := counters[winnerIdx].ack
	winner := counters[winnerIdx]
	if vots < cLevel(Quorum, st.Len()) {
		return nil, fmt.Errorf("no majority found")
	}
	if counters[winnerIdx].UTime == 0 {
		return nil, fmt.Errorf("nil object")
	}

	for _, c := range counters {
		if winner.UTime != c.UTime {
			if c.o != nil {
				wName := counters[winnerIdx].sender
				fmt.Printf("repair: receiver:%s winner:%s winnerTime %d receiverTime %d\n", c.sender, wName, winner.UTime, c.UTime)
				return winner.o, nil
				// overwrite(ctx, c.sender, winner)
			}
		}
	}
	return nil, fmt.Errorf("TODO")
}
