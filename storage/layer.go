package storage

import (
	"sync"
)

// layer is a storage primitive with optional passthrough
// to underlying layers.
//
// Layer stores values that were modified locally and
// uses recursive get() calls if the key wasn't found locally.
//
// When asked for commit(), the layer dumps its contents to the level
// under it and calls commit(). Commit wave recurses to the root level.
//
// When asked for rollback(), the parent layer is returned - and local changes are
// forgotten.
type layer struct {
	// parentLayer is this layer's parent - either
	// parent transaction or root layer.
	parentLayer *layer
	// data stores the values.
	data map[string]*valueState
	// valueCache keeps count for each unique value in the layer
	// and caches the counts coming from the underlying layers.
	valueCache map[string]uint64

	// isClosed is true if this layer was committed or rolled back.
	isClosed bool

	mu sync.Mutex
}

func newLayer() *layer {
	return &layer{
		data:       map[string]*valueState{},
		valueCache: map[string]uint64{},
	}
}

func (t *layer) set(key string, value valueState) {
	var isLocal bool
	value.Prev, isLocal = t.getIsLocal(key)

	// Crop unneeded leaves, save memory
	// 3 -> 2 -> 1 becomes 3 -> 1
	// Don't cross the layer's boundaries
	if isLocal && value.Prev != nil && value.Prev.Prev != nil {
		value.Prev = value.Prev.Prev
	}
	t.data[key] = &value
	t.refreshCacheForValue(value)
}

func (t *layer) unset(key string) {
	newValue := valueState{Data: "", Prev: t.get(key), Deleted: true}
	t.data[key] = &newValue
	t.refreshCacheForValue(newValue)
}

// get returns the value by its key.
func (t *layer) get(key string) *valueState {
	ret, _ := t.getIsLocal(key)
	return ret
}

// getIsLocal returns a valueState for the key.
// Second param is true if the value was found locally.
func (t *layer) getIsLocal(key string) (*valueState, bool) {
	// Try to return this layer's data
	ret := t.data[key]
	if ret != nil {
		return ret, true
	}

	// if no underlying layer - return
	if t.parentLayer == nil {
		return nil, false
	}

	// recurse to underlying
	return t.parentLayer.get(key), false
}

func (t *layer) numEqualTo(value string) (ret uint64) {
	// Try local storage
	val, ok := t.valueCache[value]
	if ok {
		return val
	}

	// Nowhere to recurse - return 0
	if t.parentLayer == nil {
		return 0
	}

	// Try to recurse to parentLayer
	retCount := t.parentLayer.numEqualTo(value)
	// And cache it
	t.valueCache[value] = retCount
	return retCount
}

// refreshCacheForValue actualizes the valueCache for changed values.
func (t *layer) refreshCacheForValue(value valueState) {
	// Initiate the count for current and previous values
	// from underlying layers first if exists
	if t.parentLayer != nil {
		// Init current value
		if _, ok := t.valueCache[value.Data]; !ok {
			t.valueCache[value.Data] = t.parentLayer.numEqualTo(value.Data)
		}
		// Init previous value
		if value.Prev != nil {
			if _, ok := t.valueCache[value.Prev.Data]; !ok {
				t.valueCache[value.Prev.Data] = t.parentLayer.numEqualTo(value.Prev.Data)
			}
		}
	}

	// Decrement previous value's count
	if value.Prev != nil && !value.Prev.Deleted {
		t.valueCache[value.Prev.Data]--
	}

	if !value.Deleted {
		t.valueCache[value.Data]++
	}
}
