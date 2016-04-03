package main

import (
	"github.com/ansel1/merry"
	"github.com/utrack/go-simple-memdb/storage"
)

// StorageSession handles requests for a session
// and returns output strings.
type StorageSession struct {
	stor storage.DB
}

// NewSession creates and returns new StorageSession.
func NewSession(stor storage.DB) *StorageSession {
	return &StorageSession{stor: stor}
}

// Get returns the variable's value by its key.
// Returns NULL if not found or error's text
// on unexpected error.
func (i *StorageSession) Get(key string) string {
	ret, err := i.stor.Get(key)
	if merry.Is(err, storage.ErrNotFound) {
		return "NULL"
	}
	if err != nil {
		return err.Error()
	}
	return ret
}

// Set sets the variable's value by its key.
func (i *StorageSession) Set(key, value string) {
	i.stor.Set(key, value)
}

// Unset deletes the variable by its key.
func (i *StorageSession) Unset(key string) {
	i.stor.Unset(key)
}

// NumEqualsTo returns variables' count by their value.
func (i *StorageSession) NumEqualsTo(val string) uint64 {
	return i.stor.NumEqualTo(val)
}

// Tx creates and enters new transaction.
func (i *StorageSession) Tx() string {
	i.stor = i.stor.Tx()
	return ""
}

// Commit commits current transaction in progress.
// Returns nothing on success, error on unexpected error,
// or NO TRANSACTION if not in transaction.
func (i *StorageSession) Commit() string {
	var err error
	i.stor, err = i.stor.Commit()
	if merry.Is(err, storage.ErrNoTransaction) {
		return "NO TRANSACTION"
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

// Rollback rolls back current transaction in progress (if exists).
// Returns nothing on success, error on unexpected error,
// or NO TRANSACTION if not in transaction.
func (i *StorageSession) Rollback() string {
	var err error
	i.stor, err = i.stor.Rollback()
	if merry.Is(err, storage.ErrNoTransaction) {
		return "NO TRANSACTION"
	}
	if err != nil {
		return err.Error()
	}
	return ""
}
