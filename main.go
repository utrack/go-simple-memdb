package main

import (
	"github.com/utrack/go-simple-memdb/protocol"
	"github.com/utrack/go-simple-memdb/storage"
	"os"
)

func main() {
	db := storage.New()
	// Create a protocol socket and link it to stdin/stdout
	sock := protocol.NewSocket(db)
	sock.Process(os.Stdin, os.Stdout)
}
