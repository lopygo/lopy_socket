package filter

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResult_GetDataBuffer(t *testing.T) {

}

func TestResult_GetDataLength(t *testing.T) {
}

func TestResult_SetDataBuffer(t *testing.T) {
}

func TestResult_SetPackageBuffer(t *testing.T) {
}

func TestResult_GetPackageBuffer(t *testing.T) {
}

func TestResult_GetPackageLength(t *testing.T) {
}

func TestResult_Assign(t *testing.T) {
	Convey("assign", t, func() {
		res := new(Result)
		err := res.Assign([]byte{0})
		So(err, ShouldBeError)
	})
}
