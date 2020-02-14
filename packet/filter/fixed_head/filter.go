// 头部格式固定并且包含内容长度的协议
//
// 这个目前只支持定义长度为int32的大端和小端，其它的可以参照自己写
package fixed_head

import (
	"bytes"
	"encoding/binary"
	"errors"
	"lopy_socket/packet/filter"
)


// order type
type OrderType int
const(
	// 小端，如 [03,00,00,00] 表示3
	OrderTypeLittleEndian OrderType = 0

	// 大端，如 [00,00,00,12] 表示18
	OrderTypeBigEndian    OrderType = 1
)

// new一个filter，默认为大端
func NewFilter(lengthOffset int,bodyOffset int) *Filter {
	fil := new(Filter)
	fil.lengthOffset = lengthOffset
	fil.bodyOffset = bodyOffset
	fil.orderType  = OrderTypeBigEndian
	return fil
}

// new一个小端的filter
func NewFilterLittleEndian(lengthOffset int,bodyOffset int) *Filter {
	fil := new(Filter)
	fil.lengthOffset = lengthOffset
	fil.bodyOffset = bodyOffset
	fil.orderType  = OrderTypeLittleEndian
	return fil
}

// Sentence 1
//
// Sentence 2
type Filter struct {


	// 长度值在包头的第几个字节
	lengthOffset int

	// body从哪里开始算，可能指包括一个完整包的长度，也可能只包含body的长度
	// 以此包为例 [request name(4)][length (4)][body (n)]
	// 如果length的长度是指完整包的长度，则 bodyOffset设为 0；如果length只为body的长度，则bodyOffset设为8
	bodyOffset int

	// 目前只支持大端小端
	orderType OrderType
}


func (p *Filter) GetFilterResult() (filter.IFilterResult, error) {
	result, err := NewResult(p)
	return result, err
}


func (p *Filter) Filter(buffer []byte) (filter.IFilterResult, error) {


	// 枪柄 orderType， 这是属于异常了
	err :=p.checkOrderType()
	if err != nil {
		return nil,err
	}

	// 取 result ，这个如果错了，也属于异常
	result, err := p.GetFilterResult()
	if err != nil {
		return nil, err
	}

	// 对buffer 做一些基本检查，如果错误，属于小问题
	err =  p.checkBuffer(buffer)
	if err != nil {
		// 不影响后续，所以两个nil
		return nil, nil
	}

	// 取最后一个byte的位置
	endOffset,err := p.getEndOffset(buffer)
	if err != nil {
		return nil, nil
	}

	if endOffset < p.lengthOffset+4 {
		return nil, errors.New("endOffset can not lt lengthOffset")
	}

	if endOffset < p.bodyOffset {
		return nil, errors.New("endOffset can not lt bodyOffset")
	}

	bufferLength := len(buffer)
	if bufferLength < endOffset {
		return nil,nil
	}

	packageBuffer := buffer[0 : endOffset]

	err = result.Assign(packageBuffer)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (p *Filter) checkBuffer(buffer []byte) error {
	bufferLen := len(buffer)

	// 判断长度是否达到lengthOffset要求，长度不够直接跳过
	if bufferLen < p.lengthOffset + 4 {
		return errors.New("buffer length is shorter than lengthOffset")
	}

	// 判断长度是否达到bodyOffset要求，长度不够直接跳过
	if bufferLen < p.bodyOffset {
		return errors.New("buffer length is shorter than bodyOffset")
	}

	return nil
}


func (p *Filter) checkOrderType() error {
	switch p.orderType {
	case OrderTypeBigEndian:
		return nil
	case OrderTypeLittleEndian:
		return nil
	default:
		return errors.New("orderType is not exists")
	}
}

func (p *Filter) getEndOffset(buffer []byte) (int,error) {

	var orderType binary.ByteOrder
	switch p.orderType {
	case OrderTypeBigEndian:
		orderType = binary.BigEndian
	case OrderTypeLittleEndian:
		orderType  = binary.LittleEndian
	default:
		return 0,errors.New("orderType is not exists")
	}

	err := p.checkBuffer(buffer)
	if err  != nil {
		return 0,err
	}
	lengthBuffer := new(bytes.Buffer)
	lengthBuffer.Write(buffer[p.lengthOffset:p.lengthOffset + 4])
	var bodyLength int32
	err= binary.Read(lengthBuffer,orderType,&bodyLength)
	if err != nil {
		return 0,err
	}

	endOffset := int(bodyLength) + p.bodyOffset
	return endOffset,nil
}
