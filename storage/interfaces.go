package storage

// Reader is able to retrieve values from the storage.
type Reader interface {
	// Get returns the variable's value by its key.
	// ErrNotFound is returned when the variable was not found.
	Get(key string) (string, error)
	// NumEqualTo returns the number of variables that are currently set to the passed value.
	NumEqualTo(key string) uint64
}

// Writer is able to write values to the storage.
type Writer interface {
	// Set sets the variable's value by its key.
	Set(key string, value string)
	// Unset removes the variable by its key.
	Unset(key string)
}

// ReadWriter is able to read and modify values.
type ReadWriter interface {
	Reader
	Writer
}

// DB is an instance of a database, or transaction over the DB.
// It is able to read and write values and create child transactions.
type DB interface {
	ReadWriter
	// Tx creates a transaction over the database or current transaction.
	Tx() DB
	// Commit commits the whole transaction tree, returning database's root.
	Commit() (DB, error)
	// Rollback cancels the current transaction, returning parent tx (or database's root).
	Rollback() (DB, error)
}
