package fixed_size

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewResult(t *testing.T) {
}

func TestResult_Assign(t *testing.T) {
	Convey("get result from filter",t, func() {

		iFilter,_ := NewFilter(6)

		result,err := iFilter.GetFilterResult()
		So(err,ShouldBeNil)

		Convey("赋一个nil", func() {
			err := result.Assign(nil)
			So(err,ShouldBeError)
		})

		Convey("赋一个比较短的值", func() {
			err := result.Assign([]byte{0})
			So(err,ShouldBeError)
		})

		Convey("赋一个大小相同的值", func() {
			buf := []byte{2,3,4,5,4,6}
			err := result.Assign(buf)
			So(err,ShouldBeNil)
			So(result.GetPackageBuffer(),ShouldResemble,buf)
			So(result.GetDataBuffer(),ShouldResemble,buf)
		})

		Convey("赋一个更长的值", func() {
			err := result.Assign([]byte{1,5,4,8,2,3,10,20,30,40,50,60,70})
			So(err,ShouldBeError)
		})
	})

}