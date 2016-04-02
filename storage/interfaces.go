package storage

// ReadWriter is able to retrieve values and write them back.
type ReadWriter interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Unset(key string) error
	NumEqualTo(key string) (uint64, error)
}
