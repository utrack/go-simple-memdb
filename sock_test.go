package main

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/utrack/go-simple-memdb/storage"
	"testing"
)

func TestSocket(t *testing.T) {
	Convey("With storage and socket", t, func() {
		stor := storage.New()

		sock := NewSocket(stor)

		bufIn := bytes.NewBuffer([]byte{})
		bufOut := bytes.NewBuffer([]byte{})

		Convey("Basic test", func() {
			_, _ = bufIn.WriteString(`
SET ex 10
GET ex
UNSET ex
GET ex
END
`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `UNKNOWN COMMAND

10

NULL
`)
		})

		Convey("Basic test NUMEQUALTO", func() {
			_, _ = bufIn.WriteString(`SET a 10
SET b 10
NUMEQUALTO 10
NUMEQUALTO 20
SET b 30
NUMEQUALTO 10
END`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `

2
0

1
`)
		})

		Convey("Basic tx test", func() {
			_, _ = bufIn.WriteString(`BEGIN
SET a 10
GET a
BEGIN
SET a 20
GET a
ROLLBACK
GET a
ROLLBACK
GET a
END`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `

10


20

10

NULL
`)
		})

		Convey("Tx test 2 - Double-nested Tx, ROLLBACK on no tx", func() {
			_, _ = bufIn.WriteString(`BEGIN
SET a 30
BEGIN
SET a 40
COMMIT
GET a
ROLLBACK
END`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `




40
NO TRANSACTION
`)
		})
		Convey("Tx test 3", func() {
			_, _ = bufIn.WriteString(`SET a 50
BEGIN
GET a
SET a 60
BEGIN
UNSET a
GET a
ROLLBACK
GET a
COMMIT
GET a
END`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `

50



NULL

60

60
`)
		})
		Convey("Tx test 4 - NUMEQUALTO", func() {
			_, _ = bufIn.WriteString(`SET a 10
BEGIN
NUMEQUALTO 10
BEGIN
UNSET a
NUMEQUALTO 10
ROLLBACK
NUMEQUALTO 10
COMMIT
END`)
			sock.Process(bufIn, bufOut)
			ret := bufOut.String()
			So(ret, ShouldEqual, `

1


0

1

`)
		})
	})
}
