package packet

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBlockCopy(t *testing.T) {
	Convey("测试正常情况", t, func() {
		slice1 := []byte{1,2,3,4,5,0,0,0,0,0}
		slice2 := []byte{10,20,30,40,50}

		err := BlockCopy(slice2,0,slice1,5,len(slice2))
		So(err, ShouldBeNil)
		So(len(slice1),ShouldEqual,10)
		So(slice1[5], ShouldEqual, 10)
		So(slice1[6], ShouldEqual, 20)
		So(slice1[7], ShouldEqual, 30)
		So(slice1[8], ShouldEqual, 40)
		So(slice1[9], ShouldEqual, 50)

	})

	Convey("测试索引溢出",t, func() {
		slice1 := []byte{1,2,3}
		slice2 := []byte{4,5,6}

		err1 := BlockCopy(slice1,0,slice2,0,4)
		So(err1,ShouldBeError)

		err2 := BlockCopy(slice1,4,slice2,0,0)
		So(err2,ShouldBeError)

		err3 := BlockCopy(slice1,0,slice2,4,0)
		So(err3,ShouldBeError)

	})
}
