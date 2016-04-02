package storage

import (
	"github.com/ansel1/merry"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBasicStorage(t *testing.T) {

	Convey("With layer", t, func() {
		l := newLayer()
		Convey("NUMEQUALTO for unknown value", func() {
			So(l.numEqualTo(RandString(8)), ShouldEqual, uint64(0))
		})
		Convey("Get should return nil", func() {
			ret := l.get(RandString(16))
			So(ret, ShouldBeNil)
		})
		Convey("With value", func() {
			valKey := RandString(16)
			value := ValueState{Data: RandString(64)}

			l.set(valKey, value)
			Convey("Should retrieve successfully", func() {
				got := l.get(valKey)
				So(got, ShouldResemble, &value)
			})

			Convey("Should count numEqualTo", func() {
				got := l.numEqualTo(value.Data)
				So(got, ShouldEqual, uint64(1))

				Convey("Should increment numEqualTo for same values", func() {
					valKey := RandString(16)
					l.set(valKey, value)

					got := l.numEqualTo(value.Data)
					So(got, ShouldEqual, uint64(2))

				})

				Convey("Should decrement on unset", func() {
					l.unset(valKey)

					got := l.numEqualTo(value.Data)
					So(got, ShouldEqual, uint64(0))
				})
			})

			Convey("With another value", func() {
				vk2 := RandString(16)
				v2 := ValueState{Data: RandString(256)}

				l.set(vk2, v2)

				Convey("Values should not collide", func() {
					So(l.get(valKey), ShouldResemble, &value)
					So(l.get(vk2), ShouldResemble, &v2)
				})

				Convey("Deletion shouldn't wipe the storage", func() {
					l.unset(vk2)
					So(l.get(valKey), ShouldResemble, &value)
				})
			})

			Convey("On deletion", func() {
				l.unset(valKey)
				Convey("Get should return linked value with Deleted = true", func() {
					got := l.get(valKey)
					So(got.Deleted, ShouldBeTrue)
					So(got.Prev, ShouldResemble, &value)
				})
			})

			Convey("On value's modification", func() {
				value2 := value
				value2.Data = RandString(256)

				l.set(valKey, value2)

				Convey("Should return new value linked to old one", func() {
					got := l.get(valKey)
					So(got.Data, ShouldEqual, value2.Data)
					So(got.Prev, ShouldResemble, &value)
				})

				Convey("On another mod", func() {
					value2 := value2
					value2.Data = RandString(512)
					l.set(valKey, value2)

					Convey("Should remove first mod from the link", func() {
						got := l.get(valKey)

						So(got.Data, ShouldEqual, value2.Data)
						So(got.Prev, ShouldResemble, &value)
					})
				})

				Convey("Should modify numEqualTo for mods", func() {
					So(l.numEqualTo(value.Data), ShouldEqual, uint64(0))
					So(l.numEqualTo(value2.Data), ShouldEqual, uint64(1))
				})
			})
		})
	})

}

func TestTransactions(t *testing.T) {
	Convey("With layer and fields", t, func() {
		l := newLayer()

		baseValues := map[string]ValueState{}
		var key string
		var value ValueState
		for i := 0; i < 5; i++ {
			key = RandString(16)
			value = ValueState{Data: RandString(64)}
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
					lGot, err := tx.commit(false)
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
						lGot, err := tx.commit(false)
						So(err, ShouldNotBeNil)
						So(merry.Is(err, ErrTxConflict), ShouldEqual, true)
						So(lGot, ShouldResemble, l)
					})
				})

				Convey("On non-conflicting change to the base", func() {
					l.set(RandString(32), ValueState{Data: RandString(64)})
					Convey("Commit should succeed", func() {
						_, err := tx.commit(false)
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
					lGot, err := tx2.commit(false)
					So(err, ShouldBeNil)
					So(lGot, ShouldResemble, l)
				})
			})
		})
	})

}
