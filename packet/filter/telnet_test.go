package filter

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestTelnetFilterResult_SetPackageBuffer(t *testing.T) {
	buf := []byte{1,2,3,4,5,0x0d,0x0a}
	Convey("telnet filter test",t, func() {
		a := TelnetFilterResult{}
		err := a.SetPackageBuffer(buf)
		So(err,ShouldBeNil)
		So(a.PackageBuffer(),ShouldResemble,buf)
		So(a.PackageLength(),ShouldEqual, len(buf))
		So(a.DataBuffer(),ShouldResemble,[]byte{1,2,3,4,5})
		So(a.DataLength(),ShouldEqual, len(buf) - 2)
	})

	Convey("telnet filter test",t, func() {
		a := TelnetFilterResult{}
		err := a.SetPackageBuffer([]byte{1,2,3,4,5,0x0d,0x01})
		So(err,ShouldBeError)
	})

}

func TestTelnetFilterResult_PackageBuffer(t *testing.T) {
}

func TestTelnetFilterResult_PackageLength(t *testing.T) {
}

func TestTelnetFilterResult_DataBuffer(t *testing.T) {
}

func TestTelnetFilterResult_DataLength(t *testing.T) {
}

func TestTelnetFilter_GetFilterResult(t *testing.T) {

}

func TestTelnetFilter_Filter(t *testing.T) {

	type testItem struct {
		input []byte
		isError bool
		expectPackage []byte
		expectData []byte
		expectPackageLength uint
		expectDataLength uint
	}

	arr := []testItem{
		{
			input:  []byte{1,4,6,0x0d,0x0a},
			isError: false,
			expectPackage: []byte{1,4,6,0x0d,0x0a},
			expectPackageLength: 5,
			expectData: []byte{1,4,6},
			expectDataLength: 3,
		},
		{
			input:  []byte{1,4,6,0x0d,0x0a,3,4,5},
			isError: false,
			expectPackage: []byte{1,4,6,0x0d,0x0a},
			expectPackageLength: 5,
			expectData: []byte{1,4,6},
			expectDataLength: 3,
		},

		{
			input:  []byte{1,4,6,0x0d,0x0a,3,4,5,0x0d,0x0a},
			isError: false,
			expectPackage: []byte{1,4,6,0x0d,0x0a},
			expectPackageLength: 5,
			expectData: []byte{1,4,6},
			expectDataLength: 3,
		},
	}



	for k,item := range arr {
		a := TelnetFilter{}
		res,err := a.Filter(item.input)

		var stringBuilder bytes.Buffer
		stringBuilder.WriteString("测试 telnet filter 第[")
		stringBuilder.WriteString(strconv.Itoa(k + 1))
		stringBuilder.WriteString("]次")

		Convey(stringBuilder.String(),t, func() {
			if item.isError	{
				So(err,ShouldBeError)
			}else{
				So(err,ShouldBeNil)
				So(res.DataLength(),ShouldEqual,item.expectDataLength)
				So(res.PackageLength(),ShouldEqual,item.expectPackageLength)
				So(res.PackageBuffer(),ShouldResemble,item.expectPackage)
				So(res.DataBuffer(),ShouldResemble,item.expectData)
			}
		})
	}


}