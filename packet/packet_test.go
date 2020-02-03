package packet

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInsertBuffer(t *testing.T) {
	pa := Packet{bufferZone: make([]byte,12), dataWritePosition:0}

	Convey("第一次插入数据",t, func() {
		err := pa.insertBuffer([]byte{1,2,3})
		So(err,ShouldBeNil)
		So(pa.bufferZone,ShouldResemble,[]byte{1,2,3,0,0,0,0,0,0,0,0,0})
		So(pa.dataWritePosition,ShouldEqual,3)
	})

	Convey("第二次插入数据",t, func() {
		err := pa.insertBuffer([]byte{5,6,7})
		So(err,ShouldBeNil)
		So(pa.bufferZone,ShouldResemble,[]byte{1,2,3,5,6,7,0,0,0,0,0,0})
		So(pa.dataWritePosition,ShouldEqual,6)
	})

	Convey("第三次插入数据",t, func() {
		err := pa.insertBuffer([]byte{10,20,30,40,50,60,70})
		So(err,ShouldBeNil)
		So(pa.bufferZone,ShouldResemble,[]byte{70,2,3,5,6,7,10,20,30,40,50,60})
		So(pa.dataWritePosition,ShouldEqual,1)
	})

	Convey("第四次插入数据",t, func() {
		err := pa.insertBuffer([]byte{1,2,3,4,5,6,7,8,9,10,11,12})
		So(err,ShouldBeNil)
		So(pa.bufferZone,ShouldResemble,[]byte{12,1,2,3,4,5,6,7,8,9,10,11})
		So(pa.dataWritePosition,ShouldEqual,1)
	})

	Convey("第五次插入数据",t, func() {
		err := pa.insertBuffer([]byte{0,1,2,3,4,5,6,7,8,9,10,11,12})
		So(err,ShouldBeError)
		So(pa.dataWritePosition,ShouldEqual,0)
		So(pa.dataCurrentLength,ShouldEqual,0)
	})

	Convey("第六次插入数据",t, func() {
		err := pa.insertBuffer([]byte{11,22,33})
		So(err,ShouldBeNil)
		So(pa.dataWritePosition,ShouldEqual,3)
		So(pa.bufferZone[0],ShouldEqual,11)
		So(pa.bufferZone[1],ShouldEqual,22)
		So(pa.bufferZone[2],ShouldEqual,33)
	})
}