package packet

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestNewOptionDefault(t *testing.T) {
	option := NewOptionDefault()

	Convey("测试默认option",t, func() {
		So(option.Length,ShouldEqual,4096)
	})

	Convey("测试默认option",t, func() {
		So(option.DataMaxLength,ShouldEqual,512)
	})

	Convey("测试默认情况下，check一定是正确的",t, func() {
		So(option.Check(),ShouldBeNil)
	})
}

func TestNew(t *testing.T) {


	Convey("数据最大值不能超过缓冲区大小",t, func() {
		_, err := NewOption(100,101)
		So(err,ShouldBeError)
	})

	Convey("这个测什么",t, func() {
		option, err := NewOption(100,50)
		So(err,ShouldBeNil)
		So(option.Length,ShouldEqual,100)
		So(option.Check(),ShouldBeNil)
	})
}