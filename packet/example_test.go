package packet

import (
	"fmt"
	"lopy_socket/packet/filter"
	"lopy_socket/packet/filter/fixed_head"
)

func ExampleNewPacket() {
	filterInstance := fixed_head.NewFilter(2,6)
	packetInstance := NewPacket(NewOptionDefault())
	packetInstance.SetFilter(filterInstance)

	packetInstance.OnData(func(dataResult filter.IFilterResult) {
		fmt.Println(dataResult.GetPackageBuffer())
		fmt.Println(dataResult.GetDataBuffer())
	})

	//
	fmt.Println("先试一个完整包")
	packetInstance.Put([]byte{0x23,0x23,0,0,0,2,1,2})

	fmt.Println("如果出现粘包")
	packetInstance.Put([]byte{0x24,0x24,0,0,0,2,3,4,0x25,0x25,0,0,0,3,5,6,7})

	fmt.Println("如果出现拆包")
	fmt.Println("part 1")
	packetInstance.Put([]byte{0x26,0x26,0,0,0})
	fmt.Println("part 2")
	packetInstance.Put([]byte{4,8,9,10,11})
	// 以下为打印的结果
	//先试一个完整包
	//[35 35 0 0 0 2 1 2]
	//[1 2]
	//如果出现粘包
	//[36 36 0 0 0 2 3 4]
	//[3 4]
	//[37 37 0 0 0 3 5 6 7]
	//[5 6 7]
	//如果出现拆包
	//part 1
	//part 2
	//[38 38 0 0 0 4 8 9 10 11]
	//[8 9 10 11]
}