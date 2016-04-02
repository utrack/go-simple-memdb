package storage

import (
	"fmt"
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
	data map[string]*ValueState
	// valueCache keeps count for each unique value in the layer
	// and caches the counts coming from the underlying layers.
	valueCache map[string]uint64

	mu sync.Mutex
}

func newLayer() *layer {
	return &layer{
		data:       map[string]*ValueState{},
		valueCache: map[string]uint64{},
	}
}

func (t *layer) set(key string, value ValueState) {
	value.Prev = t.get(key)

	// Crop unneeded leaves, save memory
	// 3 -> 2 -> 1 becomes 3 -> 1
	// TODO crosses tx and breaks!
	if value.Prev != nil && value.Prev.Prev != nil {
		value.Prev = value.Prev.Prev
	}
	t.data[key] = &value
	t.refreshCacheForValue(value)
}

func (t *layer) unset(key string) {
	newValue := ValueState{Data: "", Prev: t.get(key), Deleted: true}
	t.data[key] = &newValue
	t.refreshCacheForValue(newValue)
}

func (t *layer) get(key string) *ValueState {
	// Try to return this layer's data
	ret := t.data[key]
	if ret != nil {
		return ret
	}

	// if no underlying layer - return
	if t.parentLayer == nil {
		return nil
	}

	// recurse to underlying
	return t.parentLayer.get(key)
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

func (t *layer) tx() *layer {
	return &layer{
		parentLayer: t,
		data:        map[string]*ValueState{},
		valueCache:  map[string]uint64{},
	}
}

func (t *layer) commit(inRecursion bool) (ret *layer, err error) {
	// If nowhere to commit to (root layer)
	if t.parentLayer == nil {
		if inRecursion {
			// Return err if called directly
			return t, nil
		}
		return t, ErrNoTransaction.Here()
	}

	// defer recursion to parent layer's commit()
	defer func() {
		if err == nil {
			ret, err = t.parentLayer.commit(true)
		}
	}()

	// Lock the underlying layer
	t.parentLayer.mu.Lock()
	defer t.parentLayer.mu.Unlock()

	// Check for conflicts
	for key, value := range t.data {
		if t.parentLayer.get(key) != value.Prev {
			fmt.Println(value.Prev)
			fmt.Println(t.parentLayer.get(key))
			return t.parentLayer, ErrTxConflict.Here()
		}
	}

	// Copy this layer's data over and recurse
	for key, value := range t.data {
		t.parentLayer.set(key, *value)
	}

	// returns are handled in the defer
	return
}

func (t *layer) rollback() (*layer, error) {
	if t.parentLayer == nil {
		return t, ErrNoTransaction.Here()
	}
	return t.parentLayer, nil
}

// refreshCacheForValue actualizes the valueCache for changed values.
func (t *layer) refreshCacheForValue(value ValueState) {
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
