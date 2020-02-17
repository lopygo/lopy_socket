package terminator

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/lopygo/lopy_socket/packet/filter"
	"testing"
)

func TestNewFilter(t *testing.T) {
	Convey("TestNewFilter",t, func() {
		Convey("end buffer is empty", func() {
			fil,err := NewFilter([]byte{})
			So(fil,ShouldBeNil)
			So(err,ShouldNotBeNil)
		})

		Convey("end buffer is not empty", func() {

			example := []byte{4}
			fil,err := NewFilter(example)

			So(fil,ShouldImplement,(*filter.IFilter)(nil))
			So(err,ShouldBeNil)

			// length
			So(fil.endBuffer,ShouldResemble,example)
		})
	})
}

func TestFilter_Filter(t *testing.T) {
	Convey(" when the length of end buffer is 1",t, func() {

		example := []byte{4}
		Convey("当缓冲区的数据也是一位长度，且不等", func() {
			fil,_ := NewFilter(example)
			bufZone := []byte{5}

			result,err := fil.Filter(bufZone)
			So(result,ShouldBeNil)
			So(err,ShouldBeNil)

		})

		Convey("当缓冲区的数据也是一位长度，且相等", func() {
			fil,_ := NewFilter(example)

			bufZone := []byte{4}
			result,err :=  fil.Filter(bufZone)
			So(err,ShouldBeNil)
			So(result,ShouldNotBeNil)

			//
			So(result.GetDataBuffer(),ShouldResemble,[]byte{})
			So(result.GetPackageBuffer(),ShouldResemble,bufZone)

		})
	})

	Convey(" when the length of end buffer is not 1",t, func() {
		example := []byte{5,9}
		Convey("当缓冲区的数据长度相同，数据不同", func() {
			fil,_ := NewFilter(example)
			bufZone := []byte{5,7}

			result,err := fil.Filter(bufZone)
			So(result,ShouldBeNil)
			So(err,ShouldBeNil)

		})

		Convey("当缓冲区的数据长度相同，数据相同",func() {
			fil,_ := NewFilter(example)

			bufZone := []byte{5,9}
			result,err :=  fil.Filter(bufZone)
			So(err,ShouldBeNil)
			So(result,ShouldNotBeNil)

			//
			So(result.GetDataBuffer(),ShouldResemble,[]byte{})
			So(result.GetPackageBuffer(),ShouldResemble,bufZone)

		})

		Convey("当缓冲区的数据长度不同，不包含终止符", func() {
			fil,_ := NewFilter(example)

			bufZone := []byte{5,3,9}
			result,err :=  fil.Filter(bufZone)
			So(err,ShouldBeNil)
			So(result,ShouldBeNil)

			//

		})

		Convey("当缓冲区的数据长度不同，包含终止符", func() {
			fil,_ := NewFilter(example)

			bufZone := []byte{3,5,9}
			result,err :=  fil.Filter(bufZone)
			So(err,ShouldBeNil)
			So(result,ShouldNotBeNil)

			//
			So(result.GetDataBuffer(),ShouldResemble,[]byte{3})
			So(result.GetPackageBuffer(),ShouldResemble,bufZone)

		})


		Convey("当缓冲区的数据长度不同，包含终止符，且终止符后面还有数据", func() {
			fil,_ := NewFilter(example)

			bufZone := []byte{3,3,4,4,5,9,5,6,1}
			result,err :=  fil.Filter(bufZone)
			So(err,ShouldBeNil)
			So(result,ShouldNotBeNil)

			//
			So(result.GetDataBuffer(),ShouldResemble,[]byte{3,3,4,4})
			So(result.GetPackageBuffer(),ShouldResemble,[]byte{3,3,4,4,5,9})

		})


	})
}

func TestFilter_GetFilterResult(t *testing.T) {
	Convey("TestFilter_GetFilterResult",t, func() {

		fil,_ := NewFilter([]byte{1})
		result,err := fil.GetFilterResult()
		So(err,ShouldBeNil)
		So(result,ShouldImplement,(*filter.IFilterResult)(nil))

	})
}
func TestFilter_GetTerminatorBuffer(t *testing.T) {
	Convey(" TestFilter_GetTerminatorBuffer",t, func() {

		example := []byte{1}
		fil,_ := NewFilter(example)
		result := fil.GetTerminatorBuffer()

		//So(err,ShouldBeNil)
		So(result,ShouldResemble,example)

	})
}