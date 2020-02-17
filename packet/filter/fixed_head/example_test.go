package fixed_head_test

import (
	"fmt"
	"github.com/lopygo/lopy_socket/packet/filter/fixed_head"
)

func ExampleNewFilter() {
	instance := fixed_head.NewFilter(4,0,fixed_head.NewLengthTypeDefault())
	fmt.Println(instance)
}

func ExampleFilter_Filter() {

	fmt.Println("length为整个包的长度:")
	fil :=fixed_head.NewFilter(4,0,fixed_head.NewLengthTypeDefault())
	result,_ := fil.Filter([]byte{1,3,4,5,0,0,0,10,2,3})
	fmt.Println(result.GetDataBuffer())
	//[1 3 4 5 0 0 0 10 2 3]
	fmt.Println(result.GetPackageBuffer())
	//[1 3 4 5 0 0 0 10 2 3]

	fmt.Println("length为body的长度：")
	fil2 := fixed_head.NewFilter(4,8,fixed_head.NewLengthTypeDefault())
	result2,_ := fil2.Filter([]byte{1,3,4,5,0,0,0,3,1,2,3})
	fmt.Println(result2.GetDataBuffer())
	//[1 2 3]
	fmt.Println(result2.GetPackageBuffer())
	//[1 2 3 5 0 0 0 3 1 2 3]

	// Output:
	//length为整个包的长度:
	//[1 3 4 5 0 0 0 10 2 3]
	//[1 3 4 5 0 0 0 10 2 3]
	//length为body的长度：
	//[1 2 3]
	//[1 3 4 5 0 0 0 3 1 2 3]

}
