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

package objects

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/storobj"
)

// VObject is a versioned object for detecting replication inconsistencies
type VObject struct {
	// LatestObject is to most up-to-date version of an object
	LatestObject *models.Object `json:"object,omitempty"`

	// StaleUpdateTime is the LastUpdateTimeUnix of the stale object sent to the coordinator
	StaleUpdateTime int64 `json:"updateTime,omitempty"`

	// Version is the most recent incremental version number of the object
	Version uint64 `json:"version"`
}

// vobjectMarshaler is a helper for the functions implementing encoding.BinaryMarshaler
//
// Because models.Object has an optimized custom MarshalBinary function, that is what
// we want to use when serializing, rather than json.Marshal. This is just a thin
// wrapper around the model bytes resulting from the underlying call to MarshalBinary
type vobjectMarshaler struct {
	StaleUpdateTime int64
	Version         uint64
	LatestObject    []byte
}

func (vo *VObject) MarshalBinary() ([]byte, error) {
	obj, err := vo.LatestObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal object: %w", err)
	}

	b := vobjectMarshaler{vo.StaleUpdateTime, vo.Version, obj}
	return json.Marshal(b)
}

func (vo *VObject) UnmarshalBinary(data []byte) error {
	var b vobjectMarshaler

	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	vo.StaleUpdateTime = b.StaleUpdateTime
	vo.Version = b.Version

	var obj models.Object
	err = obj.UnmarshalBinary(b.LatestObject)
	if err != nil {
		return fmt.Errorf("unmarshal object: %w", err)
	}
	vo.LatestObject = &obj

	return nil
}

// Replica represents a replicated data item
type Replica struct {
	ID      strfmt.UUID     `json:"id,omitempty"`
	Deleted bool            `json:"deleted"`
	Object  *storobj.Object `json:"object,omitempty"`
}

// robjectMarshaler is a helper for the functions implementing encoding.BinaryMarshaler
//
// Because *storobj.Object has an optimized custom MarshalBinary function, that is what
// we want to use when serializing, rather than json.Marshal. This is just a thin
// wrapper around the storobj bytes resulting from the underlying call to MarshalBinary
type robjectMarshaler struct {
	ID      strfmt.UUID
	Deleted bool
	Object  []byte
}

func (ro *Replica) MarshalBinary() ([]byte, error) {
	obj, err := ro.Object.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal object: %w", err)
	}

	b := robjectMarshaler{ro.ID, ro.Deleted, obj}
	return json.Marshal(b)
}

func (ro *Replica) UnmarshalBinary(data []byte) error {
	var b robjectMarshaler

	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	ro.ID = b.ID
	ro.Deleted = b.Deleted

	var obj storobj.Object
	err = obj.UnmarshalBinary(b.Object)
	if err != nil {
		return fmt.Errorf("unmarshal object: %w", err)
	}
	ro.Object = &obj

	return nil
}

type Replicas []Replica

func (ro Replicas) MarshalBinary() ([]byte, error) {
	ms := make([]robjectMarshaler, len(ro))

	for i, obj := range ro {
		b, err := obj.Object.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("marshal object %q: %w", obj.ID, err)
		}
		ms[i] = robjectMarshaler{obj.ID, obj.Deleted, b}
	}

	return json.Marshal(ms)
}

func (ro *Replicas) UnmarshalBinary(data []byte) error {
	var ms []robjectMarshaler

	err := json.Unmarshal(data, &ms)
	if err != nil {
		return err
	}

	reps := make(Replicas, len(ms))
	for i, m := range ms {
		var obj storobj.Object
		err = obj.UnmarshalBinary(m.Object)
		if err != nil {
			return fmt.Errorf("unmarshal object %q: %w", m.ID, err)
		}
		reps[i] = Replica{
			ID:      m.ID,
			Deleted: m.Deleted,
			Object:  &obj,
		}
	}

	*ro = reps
	return nil
}

// UpdateTime return update time if it exists and 0 otherwise
func (r Replica) UpdateTime() int64 {
	if r.Object != nil {
		return r.Object.LastUpdateTimeUnix()
	}
	return 0
}
