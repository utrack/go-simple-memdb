package storage

func (t *layer) tx() *layer {
	return &layer{
		parentLayer: t,
		data:        map[string]*valueState{},
		valueCache:  map[string]uint64{},
	}
}

// commit dumps current layer's data to the parent and recurses
// commit() back to the root.
// boolean is true if commit() was called recursively.
func (t *layer) commitRecurse(inRecursion bool) (ret *layer, err error) {
	// If nowhere to commit to (root layer)
	if t.parentLayer == nil {
		if inRecursion {
			// Return err if called directly
			return t, nil
		}
		return t, ErrNoTransaction.Here()
	}

	if t.isClosed {
		return t, ErrTxClosed.Here()
	}
	// defer recursion to parent layer's commit()
	defer func() {
		if err == nil {
			ret, err = t.parentLayer.commitRecurse(true)
		}
	}()

	// Lock the underlying layer
	t.parentLayer.mu.Lock()
	defer t.parentLayer.mu.Unlock()

	// Check for conflicts
	for key, value := range t.data {
		if gotParent := t.parentLayer.get(key); gotParent != nil && gotParent != value.Prev {
			return t.parentLayer, ErrTxConflict.Here()
		}
	}

	// Copy this layer's data over and recurse
	for key, value := range t.data {
		t.parentLayer.set(key, *value)
	}

	t.isClosed = true
	// returns are handled in the defer
	return
}

// rollback returns the parent layer.
func (t *layer) rollback() (*layer, error) {
	if t.parentLayer == nil {
		return t, ErrNoTransaction.Here()
	}

	if t.isClosed {
		return t, ErrTxClosed.Here()
	}
	t.isClosed = true
	return t.parentLayer, nil
}
