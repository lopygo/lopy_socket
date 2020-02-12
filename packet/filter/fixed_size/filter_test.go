package fixed_size

import (
	. "github.com/smartystreets/goconvey/convey"
	"lopy_socket/packet/filter"
	"testing"
)

func TestNewFilter(t *testing.T) {
	Convey("new filter",t, func() {
		fi,err:= NewFilter(6)

		Convey("没有错", func() {
			So(err,ShouldBeNil)
		})
		Convey("是IFilter的接口", func() {
			So(fi,ShouldImplement,(*filter.IFilter)(nil))
		})
		Convey("是Filter", func() {
			So(fi,ShouldHaveSameTypeAs,new(Filter))
		})
	})

	Convey("new filter 错误的情况",t, func() {
		fi,err:= NewFilter(0)

		Convey("错", func() {
			So(fi,ShouldBeNil)
			So(err,ShouldBeError)
		})

	})
}

func TestFilter_GetFilterResult(t *testing.T) {
	Convey("GetFilterResult",t, func() {
		fi,_:= NewFilter(6)

		res, err := fi.GetFilterResult()
		Convey("没有错", func() {
			So(err,ShouldBeNil)
		})
		Convey("是IFilterResult的接口", func() {
			So(res,ShouldImplement,(*filter.IFilterResult)(nil))
		})
		Convey("是Result", func() {
			So(res,ShouldHaveSameTypeAs,new(Result))
		})
	})
}

func TestFilter_GetFixedSize(t *testing.T) {
	Convey("GetFixedSize",t, func() {
		fi,_:= NewFilter(6)

		Convey("是IFilterResult的接口", func() {
			So(fi.GetFixedSize(),ShouldEqual,6)
		})
	})
}

func TestFilter_Filter(t *testing.T) {
	Convey("过滤",t, func() {
		fi,_:= NewFilter(6)

		Convey("过滤一个nil", func() {
			result, err := fi.Filter(nil)
			So(err,ShouldBeNil)
			So(result,ShouldBeNil)
		})

		Convey("过滤一个比较短的值", func() {
			result, err := fi.Filter([]byte{0})
			So(err,ShouldBeNil)
			So(result,ShouldBeNil)
		})

		Convey("过滤一个大小相同的值", func() {
			buf := []byte{2,3,4,5,4,6}
			res,err := fi.Filter(buf)
			So(err,ShouldBeNil)
			So(res.GetPackageBuffer(),ShouldResemble,buf)
			So(res.GetDataBuffer(),ShouldResemble,buf)
		})

		Convey("过滤一个更长的值", func() {
			buf := []byte{1,5,4,8,2,3,10,20,30,40,50,60,70}
			res,err := fi.Filter(buf)
			So(err,ShouldBeNil)
			So(res.GetPackageBuffer(),ShouldResemble,[]byte{1,5,4,8,2,3})
			So(res.GetDataBuffer(),ShouldResemble,[]byte{1,5,4,8,2,3})
		})
	})
}