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
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/storobj"
	"github.com/weaviate/weaviate/usecases/objects"
)

type result[T any] struct {
	data T
	err  error
}

type tuple[T any] struct {
	sender string
	UTime  int64
	o      T
	ack    int
	err    error
}

type objTuple tuple[objects.Replica]

type boolTuple tuple[bool]

func readOneExists(ch <-chan simpleResult[existReply], st rState) (bool, error) {
	counters := make([]boolTuple, 0, len(st.Hosts))
	for r := range ch {
		resp := r.Response
		if r.Err != nil {
			counters = append(counters, boolTuple{resp.sender, 0, false, 0, r.Err})
			continue
		}
		counters = append(counters, boolTuple{resp.sender, resp.UpdateTime, resp.Data, 0, nil})
		max := 0
		for i := range counters {
			if r.Err == nil && counters[i].o == resp.Data {
				counters[i].ack++
			}
			if max < counters[i].ack {
				max = counters[i].ack
			}
			if max >= st.Level {
				return counters[i].o, nil
			}
		}
	}

	var sb strings.Builder
	for i, c := range counters {
		if i != 0 {
			sb.WriteString(", ")
		}
		if c.err != nil {
			fmt.Fprintf(&sb, "%s: %s", c.sender, c.err.Error())
		} else {
			fmt.Fprintf(&sb, "%s: %t", c.sender, c.o)
		}
	}
	return false, fmt.Errorf("%w %q %s", ErrConsistencyLevel, st.CLevel, sb.String())
}

type osTuple struct {
	sender string
	data   []*storobj.Object
	acks   []int
	err    error
}

func readAll(ch <-chan simpleResult[getObjectsReply], N int, st rState) ([]*storobj.Object, error) {
	ret := make([]*storobj.Object, N)
	counters := make([]osTuple, 0, len(st.Hosts))
	var sb strings.Builder
	for r := range ch {
		resp := r.Response
		if r.Err != nil {
			fmt.Fprintf(&sb, "%s: %v ", resp.sender, r.Err)
			continue
		} else if n := len(resp.Data); n != N {
			fmt.Fprintf(&sb, "%s: number of objects %d != %d ", resp.sender, n, N)
			continue
		}
		counters = append(counters, osTuple{resp.sender, resp.Data, make([]int, N), nil})
		M := 0
		for i, x := range resp.Data {
			max := 0
			for j := range counters {
				o := counters[j].data[i]
				if compare(counters[j].data[i], x) == 0 {
					counters[j].acks[i]++
				}
				if max < counters[j].acks[i] {
					max = counters[j].acks[i]
				}
				if max >= st.Level {
					ret[i] = o
				}
			}
			if max >= st.Level {
				M++
			}
		}

		if M == N {
			return ret, nil
		}
	}

	return nil, fmt.Errorf("%w %q %s", ErrConsistencyLevel, st.CLevel, sb.String())
}

// Compare returns an integer comparing two objects based on LastUpdateTimeUnix
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func compare(a, b *storobj.Object) int {
	var (
		aLastTime int64
		bLastTime int64
	)
	if a != nil {
		aLastTime = a.LastUpdateTimeUnix()
	}
	if b != nil {
		bLastTime = b.LastUpdateTimeUnix()
	}
	if aLastTime == bLastTime {
		return 0
	}
	if aLastTime > bLastTime {
		return 1
	}
	return -1
	// todo compare versions once they are implemented on objects
}
