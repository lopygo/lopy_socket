// filter规则
//
// 即拆/粘包的的规则
package filter

import "errors"

// 接口
type IFilter interface {
	// 过滤，这里两个处理结果最好都判断一下，
	//目前有三个结果 nil,nil(这种一般忽略)  IFilterResult,nil（这种为正常） nil,error（这种应该要清空缓冲区）
	Filter(buffer []byte) (IFilterResult, error)

	// 获取数据处理实例，packet将用以下实例处理数据，解决拆粘包
	GetFilterResult() (IFilterResult, error)
}

type IFilterResult interface {
	// 包的长度
	GetPackageLength() int

	// 包的buffer，包括包头包尾等，指一个完整的包
	GetPackageBuffer() []byte

	// 注入packageBuffer，
	Assign([]byte) error

	// 数据的长度，这个其实可以不要
	GetDataLength() int

	// 数据的buffer，这个其实可以不要
	GetDataBuffer() []byte
}

// 定义一个默认的Result类
type Result struct {
	packageBuffer []byte
	dataBuffer    []byte
}

func (p *Result) GetPackageLength() int {
	return len(p.packageBuffer)
}

func (p *Result) GetPackageBuffer() []byte {
	return p.packageBuffer
}

func (p *Result) SetPackageBuffer(buffer []byte) {
	p.packageBuffer = buffer
}

func (p *Result) GetDataBuffer() []byte {
	return p.dataBuffer
}

func (p *Result) SetDataBuffer(buffer []byte) {
	p.dataBuffer = buffer
}

func (p *Result) GetDataLength() int {
	return len(p.dataBuffer)
}

func (p *Result) Assign(buffer []byte) error {

	return errors.New("this method is not empty")
}
