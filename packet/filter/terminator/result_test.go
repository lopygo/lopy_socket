package terminator

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewResult(t *testing.T) {
}

func TestResult_Assign(t *testing.T) {
	Convey("get result from filter", t, func() {

		iFilter, _ := NewFilter([]byte{0, 0})

		result, err := iFilter.GetFilterResult()
		So(err, ShouldBeNil)

		Convey("赋一个比终止符短的值", func() {
			err := result.Assign([]byte{0})
			So(err, ShouldBeError)
		})

		Convey("赋一个终止符", func() {
			err := result.Assign([]byte{0, 0})
			So(err, ShouldBeNil)
			So(result.GetPackageBuffer(), ShouldResemble, []byte{0, 0})
			So(result.GetDataBuffer(), ShouldResemble, []byte{})
		})

		Convey("赋一个比终止符长的值", func() {
			err := result.Assign([]byte{1, 5, 4, 8, 0, 0})
			So(err, ShouldBeNil)
			So(result.GetPackageBuffer(), ShouldResemble, []byte{1, 5, 4, 8, 0, 0})
			So(result.GetDataBuffer(), ShouldResemble, []byte{1, 5, 4, 8})
		})

		Convey("赋一个比终止符长的值 2", func() {
			// 理论上是不会出现这种情况的
			// 实际上的应该怎么处理，算正确还是算错误。。。
			// 个人觉得应该算错误的 应该不要最后那个0
			err := result.Assign([]byte{1, 5, 4, 8, 0, 0, 0})
			So(err, ShouldBeNil)
			So(result.GetPackageBuffer(), ShouldResemble, []byte{1, 5, 4, 8, 0, 0, 0})
			So(result.GetDataBuffer(), ShouldResemble, []byte{1, 5, 4, 8, 0})
		})

		Convey("赋一个比终止符长的值，且结尾不为终止符", func() {
			err := result.Assign([]byte{1, 5, 4, 8, 0, 0})
			So(err, ShouldBeNil)
			So(result.GetPackageBuffer(), ShouldResemble, []byte{1, 5, 4, 8, 0, 0})
			So(result.GetDataBuffer(), ShouldResemble, []byte{1, 5, 4, 8})
		})
	})

}
