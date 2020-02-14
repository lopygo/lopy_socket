package fixed_head

import "fmt"

func ExampleNewFilter() {
	instance := NewFilter(4,0)
	fmt.Println(instance)
}

func ExampleFilter_Filter() {

	fil := NewFilter(4,0)
	result,_ := fil.Filter([]byte{1,3,4,5,0,0,0,10,2,3})
	fmt.Println(result.GetDataBuffer())
	//[1 3 4 5 0 0 0 10 2 3]
	fmt.Println(result.GetPackageBuffer())
	//[1 3 4 5 0 0 0 10 2 3]


	fil2 := NewFilter(4,8)
	result2,_ := fil2.Filter([]byte{1,3,4,5,0,0,0,3,1,2,3})
	fmt.Println(result2.GetDataBuffer())
	//[1 2 3]
	fmt.Println(result2.GetPackageBuffer())
	//[1 2 3 5 0 0 0 3 1 2 3]

}
