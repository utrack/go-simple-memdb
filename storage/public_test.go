package storage_test

import (
	"fmt"
	"github.com/ansel1/merry"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/utrack/go-simple-memdb/storage"
	"testing"
)

func TestPublic(t *testing.T) {
	Convey("With layer", t, func() {
		l := storage.New()

		Convey("Commit should return ErrNotInTx", func() {
			got, err := l.Commit()
			So(err, ShouldNotBeNil)
			So(merry.Is(err, storage.ErrNoTransaction), ShouldBeTrue)
			So(got, ShouldResemble, l)
		})

		Convey("Rollback should return ErrNotInTx", func() {
			got, err := l.Rollback()
			So(err, ShouldNotBeNil)
			So(merry.Is(err, storage.ErrNoTransaction), ShouldBeTrue)
			So(got, ShouldResemble, l)
		})

		Convey("Get() should return ErrNotFound", func() {
			_, err := l.Get("rand-key")
			So(err, ShouldNotBeNil)
			So(merry.Is(err, storage.ErrNotFound), ShouldBeTrue)
		})

		Convey("With set value", func() {
			key := "SomeKey"
			value := "SomeValue"
			l.Set(key, value)

			Convey("Should return the value when set", func() {
				got, err := l.Get(key)
				So(err, ShouldBeNil)
				So(got, ShouldEqual, value)
				Convey("Should return NumEqualTo", func() {
					So(l.NumEqualTo(value), ShouldEqual, uint64(1))
				})
			})
			Convey("With deleted value", func() {
				l.Unset(key)
				Convey("Should return ErrNotFound when unset", func() {
					_, err := l.Get("rand-key")
					So(err, ShouldNotBeNil)
					So(merry.Is(err, storage.ErrNotFound), ShouldBeTrue)
				})
			})
		})

	})
}

func ExampleDB() {
	db := storage.New()

	db.Set("key", "value")
	ret, _ := db.Get("key")
	fmt.Println(ret)
	// Output: value
}

func ExampleDB_transactions() {
	db := storage.New()

	// Create transaction
	tx1 := db.Tx()

	tx1.Set("key", "value")

	// Nested transactions are supported
	tx2 := tx1.Tx()
	tx2.Set("key", "value2")

	ret, _ := tx1.Get("key")
	fmt.Println(ret)

	// Transactions are committed recursively
	_, _ = tx2.Commit()
	// tx is committed now
	ret, _ = db.Get("key")
	fmt.Println(ret)

	// Output: value
	// value2
}
