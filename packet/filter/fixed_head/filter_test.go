package fixed_head

import (
	"encoding/binary"
	"github.com/lopygo/lopy_socket/packet/filter"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetOrderTypeMap(t *testing.T) {
	Convey("length > 0", t, func() {
		theMap := getOrderTypeMap()
		So(len(theMap), ShouldBeGreaterThan, 0)
	})
}

func TestCheckOrderType(t *testing.T) {

	Convey("TestCheckOrderType", t, func() {
		So(checkOrderType(OrderType(4)), ShouldBeError)
		So(checkOrderType(OrderTypeBigEndian), ShouldBeNil)
		So(checkOrderType(OrderTypeLittleEndian), ShouldBeNil)
	})
}

func TestResolveByteOrder(t *testing.T) {
	Convey("", t, func() {
		i := OrderTypeBigEndian
		j := OrderTypeLittleEndian
		e := OrderType(7)

		Convey("hehe", func() {
			So(i, ShouldEqual, OrderType(0))
			So(j, ShouldEqual, OrderType(1))
			So(e, ShouldEqual, OrderType(7))
		})

		Convey("haha i", func() {
			resI, err := ResolveByteOrder(i)

			So(err, ShouldBeNil)
			So(resI, ShouldImplement, (*binary.ByteOrder)(nil))

		})

		Convey("haha j", func() {
			resJ, err := ResolveByteOrder(i)

			So(err, ShouldBeNil)
			So(resJ, ShouldImplement, (*binary.ByteOrder)(nil))

		})

		Convey("haha e", func() {
			_, err := ResolveByteOrder(e)

			So(err, ShouldBeError)
		})

	})
}

func TestNewLengthType(t *testing.T) {
	Convey("err", t, func() {
		_, err := NewLengthType(3, OrderTypeBigEndian)
		So(err, ShouldBeError)
	})

	Convey("hehe", t, func() {

		Convey("default", func() {
			lt, err := NewLengthType(BufferLength1, OrderTypeBigEndian)
			So(err, ShouldBeNil)

			So(lt, ShouldHaveSameTypeAs, new(LengthType))
			So(lt.bufferLength, ShouldEqual, 1)
			So(lt.orderType, ShouldEqual, OrderTypeBigEndian)
		})

		Convey("BufferLength1", func() {
			lt, err := NewLengthType(BufferLength1, OrderTypeBigEndian)
			So(err, ShouldBeNil)

			So(lt, ShouldHaveSameTypeAs, new(LengthType))
			So(lt.bufferLength, ShouldEqual, 1)
			So(lt.orderType, ShouldEqual, OrderTypeBigEndian)
		})

		Convey("BufferLength2", func() {
			lt, err := NewLengthType(BufferLength2, OrderTypeBigEndian)
			So(err, ShouldBeNil)

			So(lt, ShouldHaveSameTypeAs, new(LengthType))
			So(lt.bufferLength, ShouldEqual, 2)
			So(lt.orderType, ShouldEqual, OrderTypeBigEndian)
		})

		Convey("BufferLength4", func() {
			lt, err := NewLengthType(BufferLength4, OrderTypeBigEndian)
			So(err, ShouldBeNil)

			So(lt, ShouldHaveSameTypeAs, new(LengthType))
			So(lt.bufferLength, ShouldEqual, 4)
			So(lt.orderType, ShouldEqual, OrderTypeBigEndian)
		})
	})

}

func TestLengthType_ToInt(t *testing.T) {
	Convey("ToInt 1", t, func() {
		lt, _ := NewLengthType(BufferLength1, OrderTypeLittleEndian)

		Convey("err", func() {
			_, err := lt.ToInt([]byte{1, 3})
			So(err, ShouldBeError)
		})

		Convey("normal", func() {
			res, _ := lt.ToInt([]byte{0xff})
			So(res, ShouldEqual, 255)
		})

	})

	Convey("ToInt 2", t, func() {
		lt, _ := NewLengthType(BufferLength2, OrderTypeBigEndian)

		Convey("err", func() {
			_, err := lt.ToInt([]byte{1, 3, 4})
			So(err, ShouldBeError)
		})

		Convey("normal1", func() {
			res, _ := lt.ToInt([]byte{0, 3})
			So(res, ShouldEqual, 3)
		})

		Convey("normal2", func() {
			res, _ := lt.ToInt([]byte{1, 1})
			So(res, ShouldEqual, 257)
		})

	})

	Convey("ToInt 4", t, func() {
		lt, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)

		Convey("err", func() {
			_, err := lt.ToInt([]byte{1, 3, 4})
			So(err, ShouldBeError)
		})

		Convey("normal1", func() {
			res, _ := lt.ToInt([]byte{0, 0, 0, 4})
			So(res, ShouldEqual, 4)
		})

		Convey("normal2", func() {
			res, _ := lt.ToInt([]byte{0, 0, 1, 4})
			So(res, ShouldEqual, 260)
		})

		Convey("normal4", func() {
			res, _ := lt.ToInt([]byte{1, 0, 0, 4})
			So(res, ShouldEqual, 16777220)
		})

	})
}

func TestNewFilter(t *testing.T) {
	Convey("NewFilter", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(3, 0, lengthType)
		So(fil, ShouldNotBeNil)
		So(fil, ShouldImplement, (*filter.IFilter)(nil))
		So(fil, ShouldHaveSameTypeAs, new(Filter))

		So(fil.lengthType, ShouldNotBeNil)
		So(fil.lengthType.orderType, ShouldEqual, OrderTypeBigEndian)
	})
}

func TestNewFilterLittleEndian(t *testing.T) {
	Convey("", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(3, 0, lengthType)
		So(fil, ShouldNotBeNil)
		So(fil, ShouldImplement, (*filter.IFilter)(nil))
		So(fil, ShouldHaveSameTypeAs, new(Filter))

		So(fil.lengthType, ShouldNotBeNil)
		So(fil.lengthType.orderType, ShouldEqual, OrderTypeBigEndian)
	})
}

func TestFilter_GetFilterResult(t *testing.T) {
	Convey("GetFilterResult", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		filBig := NewFilter(3, 0, lengthType)
		lengthType2, _ := NewLengthType(BufferLength4, OrderTypeLittleEndian)
		filLittle := NewFilter(3, 0, lengthType2)

		Convey("big", func() {
			result, err := filBig.GetFilterResult()
			So(err, ShouldBeNil)
			So(result, ShouldHaveSameTypeAs, new(Result))
			So(result, ShouldImplement, (*filter.IFilterResult)(nil))
		})

		Convey("little", func() {
			result, err := filLittle.GetFilterResult()
			So(err, ShouldBeNil)
			So(result, ShouldHaveSameTypeAs, new(Result))
			So(result, ShouldImplement, (*filter.IFilterResult)(nil))
		})
	})
}

func TestFilter_Filter(t *testing.T) {
	Convey("Filter", t, func() {

		Convey("lengthOffset长度不够的情况", func() {
			lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
			fil := NewFilter(4, 0, lengthType)

			//err := fil.Filter([]byte{})
			Convey("nil", func() {
				result, err := fil.Filter(nil)
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

			Convey("size 为0", func() {
				result, err := fil.Filter([]byte{})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

			Convey("size 为小于lengthOffset", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 3})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

			Convey("size 大于length，但body还没有lengthOffset的长度", func() {
				// 这个是大错，说明程序设计错了
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 1, 6, 6, 6, 6})
				So(err, ShouldBeError)
				So(result, ShouldBeNil)
			})

		})

		Convey("body长度不够的情况，bodyOffset从0开始算", func() {
			lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
			fil := NewFilter(4, 0, lengthType)

			Convey("当前buffer没有length的长", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 9})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

			Convey("当前buffer没有length的长2", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 10, 1})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

		})

		Convey("body长度不够的情况，bodyOffset不包括前面的头部", func() {
			lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
			fil := NewFilter(4, 8, lengthType)

			Convey("但body还没有lengthOffset的长度", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 1})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

			Convey("size 大于length，但body还没有lengthOffset的长度 2", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 2, 9})
				So(err, ShouldBeNil)
				So(result, ShouldBeNil)
			})

		})

		Convey("错误的情况", func() {
			Convey("endOffset < p.lengthOffset+4", func() {

				lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
				fil := NewFilter(4, 0, lengthType)
				// 这个是大错，说明程序设计错了
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 1, 6, 6, 6, 6})
				So(err, ShouldBeError)
				So(result, ShouldBeNil)
			})

			Convey("endOffset < p.bodyOffset", func() {

				lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
				// 这种情况应该是不会发生的，应该没有人这么设计
				fil := NewFilter(4, 4, lengthType)
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 3, 6, 6, 6, 6})
				So(err, ShouldBeError)
				So(result, ShouldBeNil)
			})
		})

		Convey("正常的情况", func() {
			lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
			fil := NewFilter(4, 0, lengthType)

			Convey("", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 8})
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldHaveSameTypeAs, new(Result))
			})

			Convey("1", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 9, 1})
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldHaveSameTypeAs, new(Result))
			})

			Convey("buffer length超过length的情况", func() {
				result, err := fil.Filter([]byte{1, 3, 4, 5, 0, 0, 0, 9, 1, 2, 3, 4, 5})
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldHaveSameTypeAs, new(Result))
			})
		})

	})
}

func TestFilter_CheckBuffer(t *testing.T) {

	Convey("nil", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(2, 0, lengthType)
		err := fil.checkBuffer(nil)
		So(err, ShouldBeError)

	})

	Convey("bodyOffset", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(2, 6, lengthType)

		Convey("error", func() {
			err := fil.checkBuffer([]byte{1, 2, 3, 4})
			So(err, ShouldBeError)
		})

		Convey("right", func() {
			err := fil.checkBuffer([]byte{1, 2, 3, 4, 3, 4})
			So(err, ShouldBeNil)
		})

	})

	Convey("lengthOffset", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(4, 0, lengthType)
		Convey("error", func() {
			err := fil.checkBuffer([]byte{1})
			So(err, ShouldBeError)
		})

		Convey("error 2, 长度大于等于length的开始 但小于length的结束", func() {
			err := fil.checkBuffer([]byte{1, 2, 3, 4})
			So(err, ShouldBeError)

			err = fil.checkBuffer([]byte{1, 2, 3, 4, 5, 6, 7})
			So(err, ShouldBeError)
		})

		Convey("right", func() {
			err := fil.checkBuffer([]byte{1, 2, 3, 4, 0, 0, 0, 3})
			So(err, ShouldBeNil)
		})

	})

}
func TestFilter_CheckOrderType(t *testing.T) {
	Convey("检查orderType", t, func() {
		fil := new(Filter)
		So(fil.checkLengthType(), ShouldBeError)

		Convey("when lengthType set", func() {

			fil.lengthType = NewLengthTypeDefault()

			Convey("default", func() {
				So(fil.checkLengthType(), ShouldBeNil)
			})

			Convey("小端", func() {
				fil.lengthType.orderType = OrderTypeLittleEndian
				So(fil.checkLengthType(), ShouldBeNil)
			})

			Convey("大端", func() {
				fil.lengthType.orderType = OrderTypeBigEndian
				So(fil.checkLengthType(), ShouldBeNil)
			})

			Convey("other", func() {
				fil.lengthType.orderType = 3
				So(fil.checkLengthType(), ShouldBeError)
			})

			Convey("other 2", func() {
				fil.lengthType.orderType = OrderType(5)
				So(fil.checkLengthType(), ShouldBeError)
			})
		})

	})

}
func TestFilter_GetEndOffset(t *testing.T) {

	Convey("错误的情况", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(4, 0, lengthType)
		Convey("长度为0的情况", func() {
			_, err := fil.getEndOffset([]byte{})
			So(err, ShouldBeError)
		})

		Convey("nil的情况", func() {
			_, err := fil.getEndOffset(nil)
			So(err, ShouldBeError)
		})
		Convey("长度比较小的情况", func() {
			_, err := fil.getEndOffset([]byte{2, 5})
			So(err, ShouldBeError)
		})

		Convey("长度比较小的情况，大于length的开始，小于length的结束", func() {
			_, err := fil.getEndOffset([]byte{2, 5, 4, 6, 0, 4, 0})
			So(err, ShouldBeError)
		})

	})

	Convey("大端的", t, func() {
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeBigEndian)
		fil := NewFilter(2, 0, lengthType)

		Convey("0", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 0, 0, 0, 0})
			So(offset, ShouldEqual, 0)
		})

		Convey("1", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 0, 0, 0, 8})
			So(offset, ShouldEqual, 8)
		})

		Convey("2", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 0, 0, 1, 1})
			So(offset, ShouldEqual, 257)
		})

		Convey("3", func() {
			offset, _ := fil.getEndOffset([]byte{0, 2, 0, 2, 1, 1})
			So(offset, ShouldEqual, 131329)
		})

		Convey("4", func() {
			offset, _ := fil.getEndOffset([]byte{1, 0, 4, 2, 1, 1})
			So(offset, ShouldEqual, 67240193)
		})

		Convey("5", func() {
			offset, _ := fil.getEndOffset([]byte{1, 0, 4, 2, 1, 1, 6})
			So(offset, ShouldEqual, 67240193)
		})
	})

	Convey("小端的", t, func() {
		//fil := NewFilterLittleEndian(2,0)
		//
		lengthType, _ := NewLengthType(BufferLength4, OrderTypeLittleEndian)
		fil := NewFilter(2, 0, lengthType)

		Convey("0", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 0, 0, 0, 0})
			So(offset, ShouldEqual, 0)
		})
		Convey("1", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 8, 0, 0, 0})
			So(offset, ShouldEqual, 8)
		})

		Convey("2", func() {
			offset, _ := fil.getEndOffset([]byte{0, 0, 1, 1, 0, 0})
			So(offset, ShouldEqual, 257)
		})

		Convey("3", func() {
			offset, _ := fil.getEndOffset([]byte{0, 2, 1, 1, 2, 0})
			So(offset, ShouldEqual, 131329)
		})

		Convey("4", func() {
			offset, _ := fil.getEndOffset([]byte{1, 0, 1, 1, 2, 4})
			So(offset, ShouldEqual, 67240193)
		})

		Convey("5", func() {
			offset, _ := fil.getEndOffset([]byte{1, 0, 1, 1, 2, 4, 6})
			So(offset, ShouldEqual, 67240193)
		})
	})
}
