package main

import (
	"github.com/utrack/go-simple-memdb/storage"
	"os"
)

func main() {
	db := storage.New()
	// Create a protocol socket and link it to stdin/stdout
	sock := NewSocket(db)
	sock.Process(os.Stdin, os.Stdout)
}
