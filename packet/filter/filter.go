package filter

// 接口
type IFilter interface {
	// 过滤，这里两个处理结果最好都判断一下，
	//目前有三个结果 nil,nil(这种一般忽略)  IFilterResult,nil（这种为正常） nil,error（这种应该要清空缓冲区）
	Filter(buffer []byte) (IFilterResult,error)

	// 获取数据处理实例，packet将用以下实例处理数据，解决拆粘包
	GetFilterResult() IFilterResult
}

type IFilterResult interface {
	// 包的长度
	PackageLength() uint

	// 包的buffer，包括包头包尾等，指一个完整的包
	PackageBuffer() []byte

	// 注入packageBuffer，
	SetPackageBuffer([]byte) error

	// 数据的长度，这个其实可以不要
	DataLength() uint

	// 数据的buffer，这个其实可以不要
	DataBuffer() []byte
}
