package packet

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestDefaultOption(t *testing.T) {
	option := DefaultOption()

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

func TestGetOption(t *testing.T) {


	Convey("数据最大值不能超过缓冲区大小",t, func() {
		_, err := GetOption(100,101)
		So(err,ShouldBeError)
	})

	Convey("这个测什么",t, func() {
		option, err := GetOption(100,50)
		So(err,ShouldBeNil)
		So(option.Length,ShouldEqual,100)
		So(option.Check(),ShouldBeNil)
	})
}