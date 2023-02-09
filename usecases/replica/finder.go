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
	"sort"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate/entities/additional"
	"github.com/weaviate/weaviate/entities/search"
	"github.com/weaviate/weaviate/entities/storobj"
	"github.com/weaviate/weaviate/usecases/objects"
)

// ErrConsistencyLevel consistency level cannot be achieved
var ErrConsistencyLevel = errors.New("cannot achieve consistency level")

type (
	// senderReply represent the data received from a sender
	senderReply[T any] struct {
		sender     string // hostname of the sender
		Version    int64  // sender's current version of the object
		Data       T      // the data sent by the sender
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
func (f *Finder) GetOneV2(ctx context.Context, l ConsistencyLevel, shard string,
	id strfmt.UUID, props search.SelectProperties, additional additional.Properties,
) (*storobj.Object, error) {
	c := newReadCoordinator[findOneReply](f, shard)
	op := func(ctx context.Context, host string, fullRead bool) (findOneReply, error) {
		if fullRead {
			obj, err := f.FindObject(ctx, host, f.class, shard, id, props, additional)
			var uTime int64
			if obj != nil {
				uTime = obj.LastUpdateTimeUnix()
			}
			return findOneReply{host, -1, obj, uTime, false}, err
		} else {
			xs, err := f.DigestObjects(ctx, host, f.class, shard, []strfmt.UUID{id})
			var x RepairResponse
			if len(xs) == 1 {
				x = xs[0]
			}
			if err == nil && len(xs) != 1 {
				err = fmt.Errorf("digest read request: empty result")
			}
			return findOneReply{host, x.Version, nil, x.UpdateTime, true}, err
		}
	}
	replyCh, state, err := c.Fetch2(ctx, l, op)
	if err != nil {
		return nil, err
	}
	result := <-f.readOne(ctx, shard, id, replyCh, state)
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

func (f *Finder) readOne(ctx context.Context, shard string, id strfmt.UUID, ch <-chan simpleResult[findOneReply], st rState) <-chan result[*storobj.Object] {
	// counters tracks the number of votes for each participant
	resultCh := make(chan result[*storobj.Object], 1)
	go func() {
		defer close(resultCh)
		var (
			counters = make([]objTuple, 0, len(st.Hosts))
			// winner     = 0
			max = 0
			// resultSent = false
			contentIdx = -1
		)

		for r := range ch { // len(ch) == st.Level
			resp := r.Response
			if r.Err != nil { // a least one node is not responding
				resultCh <- result[*storobj.Object]{nil, r.Err}
				return
			}
			if !resp.DigestRead {
				contentIdx = len(counters)
			}
			counters = append(counters, objTuple{resp.sender, resp.UpdateTime, resp.Data, 0, nil})
			max = 0
			for i := range counters {
				if counters[i].UTime == resp.UpdateTime {
					counters[i].ack++
				}
				if max < counters[i].ack {
					max = counters[i].ack
					// winner = i
				}
				if max >= st.Level && contentIdx >= 0 { //} !resultSent && max >= st.Level && contentIdx >= 0 {
					// resultSent = true
					resultCh <- result[*storobj.Object]{counters[contentIdx].o, nil}
					return
				}
			}
		}
		// Just in case this however should not be possible
		if n := len(counters); n == 0 || contentIdx < 0 {
			resultCh <- result[*storobj.Object]{nil, fmt.Errorf("internal error: #responses %d index %d", n, contentIdx)}
			return
		}

		obj, err := f.repairOne(ctx, shard, id, counters, st, contentIdx)
		if err == nil {
			// if !resultSent {
			resultCh <- result[*storobj.Object]{obj, nil}
			//	}
			return
		}
		// if !resultSent {
		var sb strings.Builder
		for i, c := range counters {
			if i != 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%s:%d", c.sender, c.UTime)
		}
		// sb.WriteString(err.Error())
		resultCh <- result[*storobj.Object]{nil, fmt.Errorf("%w %q: %q %v", ErrConsistencyLevel, st.CLevel, sb.String(), err)}
		//}
	}()
	return resultCh
}

// repair one object on several nodes using last write wins strategy
func (f *Finder) repairOne(ctx context.Context, shard string, id strfmt.UUID, counters []objTuple, st rState, contentIdx int) (_ *storobj.Object, err error) {
	var (
		lastUTime int64
		winnerIdx int
	)
	for i, x := range counters {
		if x.UTime > lastUTime {
			lastUTime = x.UTime
			winnerIdx = i
		}
	}

	if lastUTime < 1 {
		return nil, fmt.Errorf("nil object")
	}

	updates := counters[contentIdx].o
	winner := counters[winnerIdx]
	if updates.LastUpdateTimeUnix() != lastUTime {
		updates, err = f.RClient.FindObject(ctx, winner.sender, f.class, shard, id,
			search.SelectProperties{}, additional.Properties{})
		if err != nil {
			return nil, fmt.Errorf("get most recent object from %s: %w", counters[winnerIdx].sender, err)
		}
	}

	isNewObject := lastUTime == updates.CreationTimeUnix()
	for _, x := range counters {
		if x.UTime == 0 && !isNewObject {
			return nil, fmt.Errorf("conflict: object might have been deleted on node %q", x.sender)
		}
	}

	// case where updateTime == createTime and all other object are nil
	// TODO: if winner object is nil we need to tell the node to delete the object.
	// The adapter/repos/db/DB.overwriteObjects nil to be adjust to account for nil objects

	for _, c := range counters {
		if c.UTime != lastUTime {
			updates := []*objects.VObject{{
				LatestObject:    &updates.Object,
				StaleUpdateTime: c.UTime,
				Version:         0, // todo set when implemented
			}}
			resp, err := f.RClient.OverwriteObjects(ctx, c.sender, f.class, shard, updates)
			if err != nil {
				// fmt.Printf("repair-1: receiver:%s winner:%s winnerTime %d receiverTime %d\n", c.sender, winner.sender, winner.UTime, c.UTime)
				return nil, fmt.Errorf("node %q could not repair object: %w", c.sender, err)
			}
			if len(resp) > 0 && resp[0].Err != "" && resp[0].UpdateTime != lastUTime {
				return nil, fmt.Errorf("object changed in the meantime on node %s: %s", c.sender, resp[0].Err)
			}
			// fmt.Printf("repair-1: receiver:%s winner:%s winnerTime %d receiverTime %d\n", c.sender, winner.sender, winner.UTime, c.UTime)
			// fmt.Printf("repair: receiver:%s winner:%s winnerTime %d receiverTime %d\n", c.sender, wName, winner.UTime, c.UTime)
			// return winner.o, nil
			// overwrite(ctx, c.sender, winner)
		}
	}
	return updates, nil
}

// batchReply represents the data returned by sender
// The returned data may result from a direct or digest read request
type batchReply struct {
	// Sender hostname of the sender
	Sender string
	// DigestRead is this reply from digest read?
	DigestRead bool
	// DirectData returned from a direct read request
	DirectData []*storobj.Object
	// DigestData returned from a digest read request
	DigestData []RepairResponse
}

// GetAll gets all objects which satisfy the giving consistency
func (f *Finder) GetAllV2(ctx context.Context, l ConsistencyLevel, shard string,
	ids []strfmt.UUID,
) ([]*storobj.Object, error) {
	c := newReadCoordinator[batchReply](f, shard)
	n := len(ids)
	op := func(ctx context.Context, host string, fullRead bool) (batchReply, error) {
		if fullRead {
			xs, err := f.RClient.MultiGetObjects(ctx, host, f.class, shard, ids)
			if m := len(xs); n != m {
				err = fmt.Errorf("direct read expected %d got %d items", n, m)
			}
			return batchReply{Sender: host, DigestRead: false, DirectData: xs}, err
		} else {
			xs, err := f.DigestObjects(ctx, host, f.class, shard, ids)
			if m := len(xs); n != m {
				err = fmt.Errorf("direct read expected %d got %d items", n, m)
			}
			return batchReply{Sender: host, DigestRead: true, DigestData: xs}, err
		}
	}
	replyCh, state, err := c.Fetch2(ctx, l, op)
	if err != nil {
		return nil, err
	}
	result := <-f.readAll(ctx, shard, ids, replyCh, state)
	return result.data, result.err
}

type vote struct {
	batchReply
	Count []int
	Err   error
}

func (r batchReply) UpdateTimeAt(idx int) int64 {
	if r.DigestData != nil {
		return r.DigestData[idx].UpdateTime
	}
	return r.DirectData[idx].LastUpdateTimeUnix()
}

type _Results result[[]*storobj.Object]

func (f *Finder) readAll(ctx context.Context, shard string, ids []strfmt.UUID, ch <-chan simpleResult[batchReply], st rState) <-chan _Results {
	resultCh := make(chan _Results, 1)

	go func() {
		defer close(resultCh)
		var (
			N = len(ids) // number of requested objects
			// votes counts number of votes per object for each node
			votes      = make([]vote, 0, len(st.Hosts))
			contentIdx = -1 // index of direct read reply
		)

		for r := range ch { // len(ch) == st.Level
			resp := r.Response
			if r.Err != nil { // a least one node is not responding
				resultCh <- _Results{nil, r.Err}
				return
			}
			if !resp.DigestRead {
				contentIdx = len(votes)
			}

			votes = append(votes, vote{resp, make([]int, N), nil})
			M := 0
			for i := 0; i < N; i++ {
				max := 0
				lastTime := resp.UpdateTimeAt(i)

				for j := range votes {
					if votes[j].UpdateTimeAt(i) == lastTime {
						votes[j].Count[i]++
					}
					if max < votes[j].Count[i] {
						max = votes[j].Count[i]
					}
				}
				if max >= st.Level {
					M++
				}
			}

			if M == N {
				resultCh <- _Results{votes[contentIdx].DirectData, nil}
				return
			}
		}
		resultCh <- f.repairAll(ctx, shard, ids, votes, st, contentIdx)
	}()

	return resultCh
}

func (f *Finder) repairAll(ctx context.Context, shard string, ids []strfmt.UUID, votes []vote, st rState, contentIdx int) _Results {
	type iTuple struct {
		S int
		O int
		T int64
	}

	result := make([]*storobj.Object, len(ids))
	lastTimes := make([]iTuple, len(ids))
	// find most recent objects
	for i, x := range votes[contentIdx].DirectData {
		lastTimes[i] = iTuple{S: contentIdx, O: i, T: x.LastUpdateTimeUnix()}
	}
	for i, vote := range votes {
		if i != contentIdx {
			for j, x := range vote.DigestData {
				if x.UpdateTime > lastTimes[j].T {
					lastTimes[j] = iTuple{S: i, O: j, T: x.UpdateTime}
				}
			}
		}
	}
	// find missing content
	ys := make([]iTuple, 0, len(ids))
	for i, p := range votes[contentIdx].DirectData {
		if contentIdx != lastTimes[i].S {
			ys = append(ys, lastTimes[i])
		} else {
			result[i] = p
		}
	}
	if len(ys) > 0 {
		sort.SliceStable(ys, func(i, j int) bool { return ys[i].S < ys[j].S })
		partitions := make([]int, 0, len(votes)-2)
		pre := ys[0].S
		for i, y := range ys {
			if y.S != pre {
				partitions = append(partitions, i)
				pre = y.S
			}
		}
		partitions = append(partitions, len(ys))
		start := 0
		for _, end := range partitions {
			receiver := votes[ys[start].O].Sender
			query := make([]strfmt.UUID, end-start)
			for j := 0; start < end; start++ {
				query[j] = ids[ys[start].O]
				j++
			}
			resp, err := f.RClient.MultiGetObjects(ctx, receiver, f.class, shard, query)
			if err != nil {
				return _Results{nil, err}
			}
			n := len(query)
			if m := len(resp); n != m {
				return _Results{nil, fmt.Errorf("try to fetch %d objects from %s but got %d", m, receiver, n)}
			}
			for i := 0; i < n; i++ {
				var cTime int64
				if resp[i] != nil {
					cTime = resp[i].LastUpdateTimeUnix()
				}
				idx := ys[start-n+i].O
				if lastTimes[idx].T < cTime {
					return _Results{nil, fmt.Errorf("object %s changed on %s", ids[idx], receiver)}
				}
				result[idx] = resp[i]
			}
		}

		// preHost := votes[ys[0].S].Sender
		// os := make([]strfmt.UUID, 0, len(ys))
		// pre := ys[0].S
		// for i, y := range ys {
		// 	if y.S != pre || i == len(ys)-1 {
		// 		os[i] = ids[y.O]
		// 		os = os[:1]
		// 	} else {
		// 		os = append(os, ids[y.O])
		// 	}
		// }
	}
	return _Results{nil, nil}
}
