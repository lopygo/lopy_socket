package packet_test

import (
	"github.com/lopygo/lopy_socket/packet"
	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/terminator"
	"testing"
)

func BenchmarkTerminator(b *testing.B) {
	opt, _ := packet.NewOption(1024, 1024)
	fi, _ := terminator.NewFilter([]byte{0x0d})
	pkg := packet.NewPacket(opt)
	pkg.SetFilter(fi)
	a := 0
	pkg.OnData(func(dataResult filter.IFilterResult) {
		a++
	})
	for i := 0; i < b.N; i++ {
		err := pkg.Put([]byte{0x1, 0x2, 0x0d})
		if err != nil {
			b.Error(err)
		}
	}
	//So(a,ShouldEqual,b.N)
	if a != b.N {

		b.Error("count invalid")
	}
}


func BenchmarkTerminator2(b *testing.B) {
	opt, _ := packet.NewOption(100, 100)
	fi, _ := terminator.NewFilter([]byte{0x0d})
	pkg := packet.NewPacket(opt)
	pkg.SetFilter(fi)
	a := 0
	pkg.OnData(func(dataResult filter.IFilterResult) {
		a++
	})
	for i := 0; i < b.N; i++ {
		err := pkg.Put([]byte{0x1, 0x2, 0x0d,0x1, 0x2,0x3, 0x0d})
		if err != nil {
			b.Error(err)
		}
	}
	//So(a,ShouldEqual,b.N)
	if a != b.N * 2 {

		b.Error("count invalid")
	}
}
