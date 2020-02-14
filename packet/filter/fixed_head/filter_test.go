package fixed_head

import (
	. "github.com/smartystreets/goconvey/convey"
	"lopy_socket/packet/filter"
	"testing"
)

func TestNewFilter(t *testing.T) {
	Convey("NewFilter",t, func() {
		fil := NewFilter(3,0)
		So(fil,ShouldNotBeNil)
		So(fil,ShouldImplement,(*filter.IFilter)(nil))
		So(fil,ShouldHaveSameTypeAs,new(Filter))
		So(fil.orderType,ShouldEqual,OrderTypeBigEndian)
	})
}

func TestNewFilterLittleEndian(t *testing.T) {
	Convey("",t, func() {
		fil := NewFilterLittleEndian(3,0)
		So(fil,ShouldNotBeNil)
		So(fil,ShouldImplement,(*filter.IFilter)(nil))
		So(fil,ShouldHaveSameTypeAs,new(Filter))
		So(fil.orderType,ShouldEqual,OrderTypeLittleEndian)
	})
}

func TestFilter_GetFilterResult(t *testing.T) {
	Convey("GetFilterResult",t, func() {
		filBig := NewFilter(3,0)
		filLittle := NewFilterLittleEndian(3,0)

		Convey("big", func() {
			result,err := filBig.GetFilterResult()
			So(err,ShouldBeNil)
			So(result,ShouldHaveSameTypeAs,new(Result))
			So(result,ShouldImplement,(*filter.IFilterResult)(nil))
		})

		Convey("little", func() {
			result,err := filLittle.GetFilterResult()
			So(err,ShouldBeNil)
			So(result,ShouldHaveSameTypeAs,new(Result))
			So(result,ShouldImplement,(*filter.IFilterResult)(nil))
		})
	})
}

func TestFilter_Filter(t *testing.T) {
	Convey("Filter",t, func() {

		Convey("lengthOffset长度不够的情况", func() {
			fil := NewFilter(4,0)

			//err := fil.Filter([]byte{})
			Convey("nil", func() {
				result,err := fil.Filter(nil)
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

			Convey("size 为0", func() {
				result,err := fil.Filter([]byte{})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

			Convey("size 为小于lengthOffset", func() {
				result,err := fil.Filter([]byte{1,3,4,5,3})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

			Convey("size 大于length，但body还没有lengthOffset的长度", func() {
				// 这个是大错，说明程序设计错了
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,1,6,6,6,6})
				So(err,ShouldBeError)
				So(result,ShouldBeNil)
			})

		})

		Convey("body长度不够的情况，bodyOffset从0开始算", func() {
			fil := NewFilter(4,0)

			Convey("当前buffer没有length的长", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,9})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

			Convey("当前buffer没有length的长2", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,10,1})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})


		})

		Convey("body长度不够的情况，bodyOffset不包括前面的头部", func() {
			fil := NewFilter(4,8)

			Convey("但body还没有lengthOffset的长度", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,1})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

			Convey("size 大于length，但body还没有lengthOffset的长度 2", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,2,9})
				So(err,ShouldBeNil)
				So(result,ShouldBeNil)
			})

		})

		Convey("错误的情况", func() {
			Convey("endOffset < p.lengthOffset+4", func() {
				fil := NewFilter(4,0)
				// 这个是大错，说明程序设计错了
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,1,6,6,6,6})
				So(err,ShouldBeError)
				So(result,ShouldBeNil)
			})

			Convey("endOffset < p.bodyOffset", func() {

				// 这种情况应该是不会发生的，应该没有人这么设计
				fil := NewFilter(4,4)
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,3,6,6,6,6})
				So(err,ShouldBeError)
				So(result,ShouldBeNil)
			})
		})

		Convey("正常的情况", func() {
			fil := NewFilter(4,0)

			Convey("", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,8})
				So(err,ShouldBeNil)
				So(result,ShouldNotBeNil)
				So(result,ShouldHaveSameTypeAs,new(Result))
			})

			Convey("1", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,9,1})
				So(err,ShouldBeNil)
				So(result,ShouldNotBeNil)
				So(result,ShouldHaveSameTypeAs,new(Result))
			})

			Convey("buffer length超过length的情况", func() {
				result,err := fil.Filter([]byte{1,3,4,5,0,0,0,9,1,2,3,4,5})
				So(err,ShouldBeNil)
				So(result,ShouldNotBeNil)
				So(result,ShouldHaveSameTypeAs,new(Result))
			})
		})

	})
}

func TestFilter_CheckBuffer(t *testing.T) {

	Convey("nil",t, func() {
		fil := NewFilter(2,0)
		err := fil.checkBuffer(nil)
		So(err,ShouldBeError)

	})

	Convey("bodyOffset",t, func() {
		fil := NewFilter(2,6)

		Convey("error", func() {
			err := fil.checkBuffer([]byte{1,2,3,4})
			So(err,ShouldBeError)
		})

		Convey("right", func() {
			err := fil.checkBuffer([]byte{1,2,3,4,3,4})
			So(err,ShouldBeNil)
		})

	})

	Convey("lengthOffset",t, func() {
		fil := NewFilter(4,0)
		Convey("error", func() {
			err := fil.checkBuffer([]byte{1})
			So(err,ShouldBeError)
		})

		Convey("error 2, 长度大于等于length的开始 但小于length的结束", func() {
			err := fil.checkBuffer([]byte{1,2,3,4})
			So(err,ShouldBeError)

			err = fil.checkBuffer([]byte{1,2,3,4,5,6,7})
			So(err,ShouldBeError)
		})

		Convey("right", func() {
			err := fil.checkBuffer([]byte{1,2,3,4,0,0,0,3})
			So(err,ShouldBeNil)
		})

	})

}
func TestFilter_CheckOrderType(t *testing.T) {
	Convey("检查orderType",t, func() {
		fil := new(Filter)

		Convey("小端", func() {
			fil.orderType  = OrderTypeLittleEndian
			So(fil.checkOrderType(),ShouldBeNil)
		})

		Convey("大端", func() {
			fil.orderType  = OrderTypeBigEndian
			So(fil.checkOrderType(),ShouldBeNil)
		})

		Convey("other", func() {
			fil.orderType  = 3
			So(fil.checkOrderType(),ShouldBeError)
		})

		Convey("other 2", func() {
			fil.orderType  = OrderType(5)
			So(fil.checkOrderType(),ShouldBeError)
		})

	})

}
func TestFilter_GetEndOffset(t *testing.T) {

	Convey("错误的情况",t, func() {
		fil := NewFilter(4,0)
		Convey("长度为0的情况", func() {
			_,err := fil.getEndOffset([]byte{})
			So(err,ShouldBeError)
		})

		Convey("nil的情况", func() {
			_,err := fil.getEndOffset(nil)
			So(err,ShouldBeError)
		})
		Convey("长度比较小的情况", func() {
			_,err := fil.getEndOffset([]byte{2,5})
			So(err,ShouldBeError)
		})

		Convey("长度比较小的情况，大于length的开始，小于length的结束", func() {
			_,err := fil.getEndOffset([]byte{2,5,4,6,0,4,0})
			So(err,ShouldBeError)
		})

	})


	Convey("大端的",t, func() {
		fil := NewFilter(2,0)

		Convey("0", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,0,0,0,0})
			So(offset,ShouldEqual,0)
		})

		Convey("1", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,0,0,0,8})
			So(offset,ShouldEqual,8)
		})

		Convey("2", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,0,0,1,1})
			So(offset,ShouldEqual,257)
		})

		Convey("3", func() {
			offset,_ := fil.getEndOffset([]byte{0,2,0,2,1,1})
			So(offset,ShouldEqual,131329)
		})

		Convey("4", func() {
			offset,_ := fil.getEndOffset([]byte{1,0,4,2,1,1})
			So(offset,ShouldEqual,67240193)
		})

		Convey("5", func() {
			offset,_ := fil.getEndOffset([]byte{1,0,4,2,1,1,6})
			So(offset,ShouldEqual,67240193)
		})
	})


	Convey("小端的",t, func() {
		fil := NewFilterLittleEndian(2,0)

		Convey("0", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,0,0,0,0})
			So(offset,ShouldEqual,0)
		})
		Convey("1", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,8,0,0,0})
			So(offset,ShouldEqual,8)
		})

		Convey("2", func() {
			offset,_ := fil.getEndOffset([]byte{0,0,1,1,0,0})
			So(offset,ShouldEqual,257)
		})

		Convey("3", func() {
			offset,_ := fil.getEndOffset([]byte{0,2,1,1,2,0})
			So(offset,ShouldEqual,131329)
		})

		Convey("4", func() {
			offset,_ := fil.getEndOffset([]byte{1,0,1,1,2,4})
			So(offset,ShouldEqual,67240193)
		})

		Convey("5", func() {
			offset,_ := fil.getEndOffset([]byte{1,0,1,1,2,4,6})
			So(offset,ShouldEqual,67240193)
		})
	})
}
