package protocol

import (
	"bufio"
	"github.com/utrack/go-simple-memdb/storage"
	"io"
	"strconv"
	"strings"
)

// DBSocket is a sock scanner that reads commands and returns their output.
type DBSocket struct {
	sess *StorageSession
}

// NewSocket returns new DBSocket that reads requests
// according to protocol specs, relays
// them to the Database via StorageSession and returns the output.
func NewSocket(db storage.DB) *DBSocket {
	return &DBSocket{sess: NewSession(db)}
}

// Process starts the IO pipe.
func (s *DBSocket) Process(rPipe io.Reader, wPipe io.Writer) {
	r := bufio.NewReader(rPipe)
	w := bufio.NewWriter(wPipe)

	var cmdRaw string
	var cmd []string
	var err error
	for {
		cmdRaw, err = r.ReadString('\n')
		cmdRaw = strings.Trim(cmdRaw, "\n ")
		if err != nil {
			return
		}

		cmd = strings.SplitN(cmdRaw, " ", 3)
		output := ""
		switch cmd[0] {
		case "END":
			return
		case "GET":
			output = s.sess.Get(cmd[1])
		case "SET":
			s.sess.Set(cmd[1], cmd[2])
		case "UNSET":
			s.sess.Unset(cmd[1])
		case "NUMEQUALTO":
			output = strconv.FormatUint(s.sess.NumEqualsTo(cmd[1]), 10)
		case "BEGIN":
			output = s.sess.Tx()
		case "COMMIT":
			output = s.sess.Commit()
		case "ROLLBACK":
			output = s.sess.Rollback()
		default:
			output = "UNKNOWN COMMAND"
		}
		_, _ = w.WriteString(output)
		_ = w.WriteByte('\n')
		_ = w.Flush()
	}
}
