package storage

// New creates new storage instance.
func New() DB {
	return newLayer()
}

// Commit implements DB interface.
func (t *layer) Commit() (DB, error) {
	return t.commitRecurse(false)
}

// Rollback implements DB interface.
func (t *layer) Rollback() (DB, error) {
	return t.rollback()
}

// Get implements Reader interface.
func (t *layer) Get(key string) (string, error) {
	ret := t.get(key)
	if ret == nil || ret.Deleted {
		return ``, ErrNotFound.Here()
	}
	return ret.Data, nil
}

// NumEqualTo implements Reader interface.
func (t *layer) NumEqualTo(value string) uint64 {
	return t.numEqualTo(value)
}

// Set implements Writer interface.
func (t *layer) Set(key, value string) {
	t.set(key, valueState{Data: value})
}

// Unset implements Writer interface.
func (t *layer) Unset(key string) {
	t.unset(key)
}

// Tx implements DB interface.
func (t *layer) Tx() DB {
	return t.tx()
}
