package storage

import (
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
			value := valueState{Data: RandString(64)}

			l.set(valKey, value)
			Convey("getIsLocal should return true", func() {
				got, isLocal := l.getIsLocal(valKey)
				So(got, ShouldResemble, &value)
				So(isLocal, ShouldBeTrue)
			})
			Convey("Should retrieve successfully", func() {
				got := l.get(valKey)
				So(got, ShouldResemble, &value)
			})

			Convey("Deletion", func() {
				l.unset(valKey)

				Convey("Should return deleted value", func() {
					got := l.get(valKey)
					So(got.Deleted, ShouldEqual, true)
					So(got.Data, ShouldEqual, "")
				})

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
				v2 := valueState{Data: RandString(256)}

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
