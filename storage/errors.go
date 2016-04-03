package storage

import (
	"github.com/ansel1/merry"
)

// ErrNotFound is returned if the key was not found.
var ErrNotFound = merry.New("Key was not found!")

// ErrNoTransaction is returned when tx-specific function
// like Commit() or Rollback() was called without
// any transactions present.
var ErrNoTransaction = merry.New("There is no transaction in progress.")

// ErrTxConflict is returned when a conflict was found
// during the Commit().
var ErrTxConflict = merry.New("Transaction conflict! Aborted.")

// ErrTxClosed is returned when trying to commit transaction
// that was committed before.
var ErrTxClosed = merry.New("Transaction was closed.")
