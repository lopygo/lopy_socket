package telnet

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"lopy_socket/packet/filter"
	"lopy_socket/packet/filter/terminator"

	//"lopy_socket/packet/filter/terminator"
	"testing"
)

func TestNew(t *testing.T) {

	Convey("测试 telnet filter",t, func() {
		a,err := NewFilter()
		Convey("检测类型", func() {
			So(err, ShouldBeNil)
			So(a,ShouldImplement,(*filter.IFilter)(nil))
			So(a,ShouldHaveSameTypeAs,new(terminator.Filter))
			//So(a,ShouldImplement,terminator.Filter{})
		})

		Convey("试一下filter结果，正常结果", func() {
			buf := []byte{2,3,7,0x0a}
			res,err := a.Filter(buf)

			So(err,ShouldBeNil)
			So(res.GetDataBuffer(),ShouldResemble,[]byte{2,3,7})
			So(res.GetPackageBuffer(),ShouldResemble,[]byte{2,3,7,0x0a})
		})

		Convey("试一下filter结果，没有终止符", func() {
			buf := []byte{2,3,7}
			res,err := a.Filter(buf)

			So(err,ShouldBeNil)
			So(res,ShouldBeNil)
		})

		Convey("试一下filter结果，后面还有数据", func() {
			buf := new(bytes.Buffer)
			buf.WriteString("hello\nworld\n")

			res,err := a.Filter(buf.Bytes())
			res.GetDataBuffer()

			dataBuf :=new(bytes.Buffer)
			dataBuf.WriteString("hello")

			packageBuf :=new(bytes.Buffer)
			packageBuf.WriteString("hello\n")

			So(err,ShouldBeNil)
			So(res.GetDataBuffer(),ShouldResemble,dataBuf.Bytes())
			So(res.GetPackageBuffer(),ShouldResemble,packageBuf.Bytes())
		})
	})

}