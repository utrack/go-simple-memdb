package storage

// ValueState represents a unique value in the storage.
type ValueState struct {
	// Data is this value's underlying data.
	Data string
	// Prev is a ptr to the previous state of this value.
	Prev *ValueState

	// Deleted is true if the value had been deleted.
	// Treat as NULL.
	Deleted bool
}
