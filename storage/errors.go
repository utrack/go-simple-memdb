package storage

import (
	"github.com/ansel1/merry"
)

// ErrNotFound is returned if the key was not found.
var ErrNotFound = merry.New("Key was not found!")

// ErrNoTransaction is returned when tx-specific function like Commit() or Rollback()
// was called without any transactions present.
var ErrNoTransaction = merry.New("There is no transaction in progress.")

// ErrTxConflict is returned when TODO
var ErrTxConflict = merry.New("Transaction conflict! Aborted.")
