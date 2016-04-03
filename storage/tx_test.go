package storage

import (
	"github.com/ansel1/merry"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTransactions(t *testing.T) {
	Convey("With layer and fields", t, func() {
		l := newLayer()

		baseValues := map[string]valueState{}
		var key string
		var value valueState
		for i := 0; i < 5; i++ {
			key = RandString(16)
			value = valueState{Data: RandString(64)}
			baseValues[key] = value
			l.set(key, value)
		}

		Convey("With child layer", func() {
			tx := l.tx()
			Convey("Should retrieve valid NUMEQUALTO", func() {
				for _, value := range baseValues {
					So(tx.numEqualTo(value.Data), ShouldEqual, uint64(1))
				}
			})

			Convey("Should not commit twice", func() {
				_, _ = tx.commitRecurse(false)
				_, err := tx.commitRecurse(false)
				So(err, ShouldNotBeNil)
				So(merry.Is(err, ErrTxClosed), ShouldBeTrue)
			})

			Convey("getIsLocal should return false", func() {
				got, isLocal := tx.getIsLocal(key)
				So(got, ShouldResemble, &value)
				So(isLocal, ShouldBeFalse)
			})

			Convey("On change for tx", func() {
				newValue := value
				newValue.Data = RandString(256)
				tx.set(key, newValue)
				Convey("Should return the value for tx", func() {
					got := tx.get(key)
					So(got.Data, ShouldEqual, newValue.Data)
				})
				Convey("Should return old value for base", func() {
					got := l.get(key)
					So(got.Data, ShouldEqual, value.Data)
				})
				Convey("Should return sane counts", func() {
					So(tx.numEqualTo(newValue.Data), ShouldEqual, uint64(1))
					So(l.numEqualTo(newValue.Data), ShouldEqual, uint64(0))
					So(tx.numEqualTo(value.Data), ShouldEqual, uint64(0))
					So(l.numEqualTo(value.Data), ShouldEqual, uint64(1))
				})

				Convey("Commit", func() {
					lGot, err := tx.commitRecurse(false)
					So(err, ShouldBeNil)
					So(lGot, ShouldResemble, l)
					Convey("Tx values should be written", func() {
						got := l.get(key)
						So(got.Data, ShouldEqual, newValue.Data)
					})
					Convey("Counts should be updated", func() {
						So(l.numEqualTo(newValue.Data), ShouldEqual, uint64(1))
						So(l.numEqualTo(value.Data), ShouldEqual, uint64(0))
					})
				})
				Convey("Rollback", func() {
					lGot, err := tx.rollback()
					So(err, ShouldBeNil)
					So(lGot, ShouldResemble, l)
					Convey("Changes should be forgotten", func() {
						got := l.get(key)
						So(got.Data, ShouldEqual, value.Data)
					})
				})

				Convey("On conflicting change to the base", func() {
					value.Data = RandString(256)
					l.set(key, value)
					Convey("Commit should fail", func() {
						lGot, err := tx.commitRecurse(false)
						So(err, ShouldNotBeNil)
						So(merry.Is(err, ErrTxConflict), ShouldEqual, true)
						So(lGot, ShouldResemble, l)
					})
				})

				Convey("On non-conflicting change to the base", func() {
					l.set(RandString(32), valueState{Data: RandString(64)})
					Convey("Commit should succeed", func() {
						_, err := tx.commitRecurse(false)
						So(err, ShouldBeNil)
					})

				})
			})

			Convey("With another child layer", func() {
				tx2 := tx.tx()

				Convey("Rollback should return parent's layer", func() {
					lGot, err := tx2.rollback()
					So(err, ShouldBeNil)
					So(lGot, ShouldResemble, tx)
				})
				Convey("Commit should return base layer", func() {
					lGot, err := tx2.commitRecurse(false)
					So(err, ShouldBeNil)
					So(lGot, ShouldResemble, l)
				})

			})
		})
	})

}
