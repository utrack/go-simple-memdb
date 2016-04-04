package protocol

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/utrack/go-simple-memdb/storage"
	"testing"
)

// testStorage is a mock for storage.DB.
// It calls respective functions for the session calls.
type testStorage struct {
	fGet   func(string) (string, error)
	fSet   func(string, string)
	fUnset func(string)

	fNumEqualTo func(string) uint64

	fTx       func() storage.DB
	fCommit   func() (storage.DB, error)
	fRollback func() (storage.DB, error)
}

func (t *testStorage) Get(key string) (string, error) {
	return t.fGet(key)
}

func (t *testStorage) Set(key, value string) {
	t.fSet(key, value)
}

func (t *testStorage) Unset(key string) {
	t.fUnset(key)
}

func (t *testStorage) NumEqualTo(value string) uint64 {
	return t.fNumEqualTo(value)
}

func (t *testStorage) Tx() storage.DB {
	return t.fTx()
}

func (t *testStorage) Commit() (storage.DB, error) {
	return t.fCommit()
}

func (t *testStorage) Rollback() (storage.DB, error) {
	return t.fRollback()
}

// We set one function at a time and test StorageSession
// by calling respective functions.
// If something unexpected was called - tests should fail
// with nullptr panic.
func TestStorageSession(t *testing.T) {

	Convey("With mock storage and session handler", t, func() {
		s := &testStorage{}
		sessHandler := NewSession(s)

		Convey("Get", func() {
			var gotKey string
			s.fGet = func(k string) (string, error) {
				gotKey = k
				return "someVal", nil
			}

			sentKey := "some-key"
			got := sessHandler.Get(sentKey)

			So(gotKey, ShouldEqual, sentKey)
			So(got, ShouldEqual, "someVal")
			Convey("Should return NULL if not found", func() {
				s.fGet = func(_ string) (string, error) {
					return "", storage.ErrNotFound.Here()
				}
				got := sessHandler.Get("k")
				So(got, ShouldEqual, "NULL")

			})
		})

		Convey("Set", func() {
			var gotKey string
			var gotVal string
			s.fSet = func(k, v string) {
				gotKey = k
				gotVal = v
				return
			}

			sentKey := "someKey"
			sentVal := "someVal"

			sessHandler.Set(sentKey, sentVal)
			So(gotKey, ShouldEqual, sentKey)
			So(gotVal, ShouldEqual, sentVal)

		})

		Convey("Unset", func() {
			var gotKey string
			s.fUnset = func(k string) {
				gotKey = k
			}
			sentKey := "some-key-sent"
			sessHandler.Unset(sentKey)
			So(gotKey, ShouldEqual, sentKey)
		})

		Convey("NumEqualTo", func() {
			var gotKey string
			s.fNumEqualTo = func(k string) uint64 {
				gotKey = k
				return 1337
			}

			sentKey := "sent-key"
			ret := sessHandler.NumEqualsTo(sentKey)
			So(ret, ShouldEqual, uint64(1337))
			So(gotKey, ShouldEqual, sentKey)
		})

		Convey("Tx should assign returned storage to stor", func() {
			sentStor := &testStorage{}
			s.fTx = func() storage.DB {
				return sentStor
			}
			got := sessHandler.Tx()
			So(got, ShouldEqual, "")
			So(sessHandler.stor, ShouldResemble, sentStor)
		})

		Convey("Commit", func() {
			sentStor := &testStorage{}
			s.fCommit = func() (storage.DB, error) {
				return sentStor, nil
			}
			got := sessHandler.Commit()
			So(got, ShouldEqual, "")
			So(sessHandler.stor, ShouldResemble, sentStor)
			Convey("Proper NO TRANSACTION", func() {
				sessHandler.stor = s
				s.fCommit = func() (storage.DB, error) {
					return sentStor, storage.ErrNoTransaction.Here()
				}
				got := sessHandler.Commit()
				So(got, ShouldEqual, "NO TRANSACTION")
				So(sessHandler.stor, ShouldEqual, sentStor)
			})
		})
		Convey("Rollback", func() {
			sentStor := &testStorage{}
			s.fRollback = func() (storage.DB, error) {
				return sentStor, nil
			}
			got := sessHandler.Rollback()
			So(got, ShouldEqual, "")
			So(sessHandler.stor, ShouldResemble, sentStor)
			Convey("Proper NO TRANSACTION", func() {
				sessHandler.stor = s
				s.fRollback = func() (storage.DB, error) {
					return sentStor, storage.ErrNoTransaction.Here()
				}
				got := sessHandler.Rollback()
				So(got, ShouldEqual, "NO TRANSACTION")
				So(sessHandler.stor, ShouldEqual, sentStor)
			})
		})
	})

}
