package storage

// Reader is able to retrieve values from the storage.
type Reader interface {
	Get(key string) (string, error)
	NumEqualTo(key string) (uint64, error)
}

// Writer is able to write values to the storage.
type Writer interface {
	Set(key string, value string) error
	Unset(key string) error
}

// ReadWriter is able to read and modify values.
type ReadWriter interface {
	Reader
	Writer
}

// DB is an instance of a database, or transaction over it.
type DB interface {
	ReadWriter
	// Tx creates a transaction over the database or current transaction.
	Tx() DB
	// Commit commits the transaction tree, returning database's root.
	Commit() (DB, error)
	// Rollback cancels the transaction, returning parent tx (or database's root).
	Rollback() (DB, error)
}
