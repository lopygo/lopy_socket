package fixed_head

import (
	. "github.com/smartystreets/goconvey/convey"
	"lopy_socket/packet/filter"
	"testing"
)


func TestNewResult(t *testing.T) {
	Convey("new 一个result",t, func() {
		iFilter := NewFilter(2,0,NewLengthTypeDefault())
		result,err := NewResult(iFilter)

		So(err,ShouldBeNil)
		So(result,ShouldNotBeNil)
		So(result,ShouldImplement,(*filter.IFilterResult)(nil))
		So(result,ShouldHaveSameTypeAs,new(Result))

	})
}

func TestResult_Assign(t *testing.T) {
	Convey("body长度为整个包",t, func() {

		iFilter := NewFilter(2,0,NewLengthTypeDefault())
		result,err := NewResult(iFilter)

		So(err,ShouldBeNil)

		Convey("赋一个nil", func() {
			err := result.Assign(nil)
			So(err,ShouldBeError)
		})

		Convey("赋一个比较短的值", func() {
			err := result.Assign([]byte{0})
			So(err,ShouldBeError)
		})


		//
		Convey("赋一个完整包的值", func() {
			buf := []byte{2,3,0,0,0,8,4,5}
			err := result.Assign(buf)
			So(err,ShouldBeNil)
			So(result.GetPackageBuffer(),ShouldResemble,buf)
			So(result.GetDataBuffer(),ShouldResemble,buf)
		})

		Convey("赋一个比完整包size大的值", func() {
			buf := []byte{2,3,0,0,0,8,4,5,1}
			err := result.Assign(buf)
			So(err,ShouldBeError)
		})

	})

	Convey("body长度不包含head",t, func() {

		iFilter := NewFilter(2,6,NewLengthTypeDefault())
		result,err := NewResult(iFilter)

		So(err,ShouldBeNil)

		Convey("赋一个nil", func() {
			err := result.Assign(nil)
			So(err,ShouldBeError)
		})

		Convey("赋一个比较短的值", func() {
			err := result.Assign([]byte{0})
			So(err,ShouldBeError)
		})


		//
		Convey("赋一个完整包的值", func() {
			buf := []byte{2,3,0,0,0,2,4,5}
			err := result.Assign(buf)
			So(err,ShouldBeNil)
			So(result.GetPackageBuffer(),ShouldResemble,buf)
			So(result.GetDataBuffer(),ShouldResemble,[]byte{4,5})
		})

	})

}