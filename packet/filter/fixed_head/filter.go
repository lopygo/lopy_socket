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
//
// 看了binary的源码，这里的OrderType应该不用导出，所以orderType应该用小写
type OrderType int
const(

	// 大端，如 [00,00,00,12] 表示18
	OrderTypeBigEndian    OrderType = 0

	// 小端，如 [03,00,00,00] 表示3
	OrderTypeLittleEndian OrderType = 1

)

// 这个是不是没有什么卵用
func getOrderTypeMap() (orderTypeMap map[OrderType]binary.ByteOrder) {
	orderTypeMap  =make(map[OrderType]binary.ByteOrder)
	orderTypeMap[OrderTypeBigEndian]  = binary.BigEndian
	orderTypeMap[OrderTypeLittleEndian]  = binary.LittleEndian

	return orderTypeMap
}

func checkOrderType(orderType OrderType) error {
	orderTypeMap  :=getOrderTypeMap()

	_,ok := orderTypeMap[orderType]
	if ok {
		return nil
	}

	return errors.New("orderType error")
}

func ResolveByteOrder(receiver OrderType) (binary.ByteOrder,error)  {
	var byteOrder binary.ByteOrder

	orderTypeMap := getOrderTypeMap()

	byteOrder,ok := orderTypeMap[receiver]
	if !ok {
		return nil,errors.New("orderType is not exists")
	}
	return byteOrder,nil
}

type BufferLength uint16
const(
	BufferLength1 = BufferLength(1)
	BufferLength2 = BufferLength(2)
	BufferLength4 = BufferLength(4)
)

func getBufferLengthMap() (bufferLengthMap map[BufferLength]BufferLength) {
	bufferLengthMap =make(map[BufferLength]BufferLength)
	bufferLengthMap[BufferLength1]  = BufferLength1
	bufferLengthMap[BufferLength2]  = BufferLength2
	bufferLengthMap[BufferLength4]  = BufferLength4

	return bufferLengthMap
}

type LengthType struct {
	bufferLength BufferLength // 1,2,4
	orderType OrderType
}

func NewLengthTypeDefault() (res *LengthType) {
	res,_ = NewLengthType(BufferLength4,OrderTypeBigEndian)
	return
}


func NewLengthType(bufferLength BufferLength,orderType OrderType) (res *LengthType,err error) {
	res = nil
	maps := getBufferLengthMap()

	//
	_,ok := maps[bufferLength]
	if !ok {
		err = errors.New("")
		return
	}

	err = checkOrderType(orderType)
	if err != nil{
		return
	}


	res = new(LengthType)
	res.bufferLength = bufferLength
	res.orderType = orderType

	return
}


func (p *LengthType) ToInt(buffer []byte) ( int,error) {
	length := len(buffer)

	// 长度判断
	if length != int(p.bufferLength) {
		return 0,errors.New("buffer length error")
	}

	orderType,err := ResolveByteOrder(p.orderType)
	if  err != nil{
		return 0,err
	}

	lengthBuffer := new(bytes.Buffer)
	lengthBuffer.Write(buffer[0: length])


	// 这个这么定义，在32位系统上应该会出问题吧
	var res int  = 0

	switch p.bufferLength {
	case BufferLength1:
		var tmpLen uint8 =0
		err = binary.Read(lengthBuffer,orderType,&tmpLen)
		res = int(tmpLen)
	case BufferLength2:
		var tmpLen uint16 =0
		err = binary.Read(lengthBuffer,orderType,&tmpLen)
		res = int(tmpLen)
	case BufferLength4:
		var tmpLen uint32 =0
		err = binary.Read(lengthBuffer,orderType,&tmpLen)
		res = int(tmpLen)

	default:
		return 0,errors.New("length error")
	}

	return res,nil
}

// new一个filter，默认为大端
func NewFilter(lengthOffset int,bodyOffset int,lengthType *LengthType) *Filter {
	fil := new(Filter)
	fil.lengthOffset = lengthOffset
	fil.bodyOffset = bodyOffset

	fil.lengthType = lengthType
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


	// length type, length 和 Order的组合
	lengthType *LengthType
}


func (p *Filter) GetFilterResult() (filter.IFilterResult, error) {
	result, err := NewResult(p)
	return result, err
}


func (p *Filter) Filter(buffer []byte) (filter.IFilterResult, error) {


	// 枪柄 orderType， 这是属于异常了
	err :=p.checkLengthType()
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


func (p *Filter) checkLengthType() error {
	if nil == p.lengthType {
		return errors.New("lengthType can not be empty")
	}

	err := checkOrderType(p.lengthType.orderType)

	if err != nil {
		return err
	}

	return nil
}

func (p *Filter) getEndOffset(buffer []byte) (int,error) {

	//var orderType binary.ByteOrder
	//switch p.orderType {
	//case OrderTypeBigEndian:
	//	orderType = binary.BigEndian
	//case OrderTypeLittleEndian:
	//	orderType  = binary.LittleEndian
	//default:
	//	return 0,errors.New("orderType is not exists")
	//}



	err := p.checkBuffer(buffer)
	if err  != nil {
		return 0,err
	}

	bodyLength,err:= p.lengthType.ToInt(buffer[p.lengthOffset:p.lengthOffset + int(p.lengthType.bufferLength)])
	//
	//lengthBuffer := new(bytes.Buffer)
	//lengthBuffer.Write(buffer[p.lengthOffset:p.lengthOffset + 4])
	//var bodyLength int32
	//err= binary.Read(lengthBuffer,orderType,&bodyLength)
	if err != nil {
		return 0,err
	}

	endOffset := bodyLength + p.bodyOffset
	return endOffset,nil
}
