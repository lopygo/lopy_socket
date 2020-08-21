// 用于解决tcp拆/粘包的一个练手项目
//
// 能方便的完成拆包
//
// 刚学go，思想还不太熟悉，先将就着用
//
// 目前用代码为同步方式，异步的(goroutine)的还不知道怎么测试，后面再考虑吧
//
// 快速开始
//	// hehe,后面再补吧
//	fmp.Println("hello world")
package packet

import (
	"errors"
	pkgBuffer "github.com/lopygo/lopy_socket/packet/buffer"
	"github.com/lopygo/lopy_socket/packet/filter"
	"log"
	"sync"
)

// 这是什么
type Packet struct {
	// 定义一个Filter类
	dataFilter filter.IFilter

	// 数据读取点，也就是数据在缓冲区的起点index
	dataReadPosition int

	// 写数据的点，也就是写之前数据在缓冲区的终点
	dataWritePosition int

	// 缓冲区
	bufferZone []byte

	// 数据包的最大长度，表示如果一个包超过了这个长度，那么将被丢弃
	dataMaxLength int

	// 回调函数，先用这种简单的方式，回头用event之类的
	onDataCallback func(data filter.IFilterResult)

	// 缓冲区的最大空间，如果要做自动扩容，那么，这个配置表示扩容后的最大空间，暂时不做自动扩容
	//bufferZoneMaxLength int

	m sync.Mutex
}

func NewPacket(option *Option) *Packet {
	packetInstance := new(Packet)
	packetInstance.dataMaxLength = option.DataMaxLength
	packetInstance.bufferZone = make([]byte, option.Length, option.Length)
	packetInstance.dataFilter = option.Filter

	return packetInstance
}

// 缓冲区长度
func (p *Packet) bufferZoneLength() int {
	return len(p.bufferZone)
}

// Deprecated: 设置过滤规则，即 粘/拆包的过滤规则
//
// 这个方法不应该这样写的，而应该在NewPacket的时候指定
//
// 当初这么写是对语言的不熟，（那会还不知道怎么实现构造方法。。。）
//
// 下一步可以考虑把filter写进option
func (p *Packet) SetFilter(filter filter.IFilter) {
	p.dataFilter = filter
}

// 获取过滤规则，即 粘/拆包的规则
func (p *Packet) GetFilter() (filter.IFilter, error) {

	if nil == p.dataFilter {
		return nil, errors.New("filter must be set")
	}
	return p.dataFilter, nil
}

// 获取可用长度，即总的缓冲区长度减当前数据长度
//
// 现在加一个，长度减1
func (p *Packet) getAvailableLen() int {
	bufLen := p.bufferZoneLength()
	dataLen := p.currentDataLength()
	if dataLen >= bufLen {
		return 0
	}

	return p.bufferZoneLength() - p.currentDataLength() - 1
}

func (p *Packet) GetAvailableLen() int {
	p.m.Lock()
	l := p.getAvailableLen()
	p.m.Unlock()
	return l
}

// 这个看怎么做，先暂时写这种简单的
func (p *Packet) OnData(callback func(dataResult filter.IFilterResult)) {
	p.onDataCallback = callback
}

// 外部写入数据
//
// 要考虑并发问题，目前没有考虑
//
// 后面先试一试锁，然后再考虑channel
func (p *Packet) Put(data []byte) error {

	// 为空，则不管
	if nil == data || len(data) == 0 {
		return nil
	}

	dataLength := len(data)
	if dataLength > p.dataMaxLength {
		msg := "数据长度限制，请跟据实际情况重新配置 dataMaxLength"
		return errors.New(msg)
	}

	//fmt.Println(p.getAvailableLen())
	//fmt.Println(p.bufferZoneLength() - p.currentDataLength())
	//if dataLength+p.currentDataLength() > p.bufferZoneLength() {
	if dataLength > p.getAvailableLen() {
		msg := "缓冲区大小不足"
		return errors.New(msg)
	}

	p.m.Lock()
	defer p.m.Unlock()

	// 插入数据
	err := p.insertBuffer(data)
	if err != nil {
		return err
	}

	// filter
	p.readByFilter()

	return nil
}

// 读出filter后的数据
func (p *Packet) readByFilter() {
	for {
		data, err := p.getCurrentData()
		if err != nil {
			break
		}
		if len(data) == 0 {
			break
		}

		// 从这里开始，的异常，要清空数据
		fil, _ := p.GetFilter()
		if fil == nil {
			// flush
			p.dataReadPosition = 0
			p.dataWritePosition = 0
			break
		}
		filterResult, err := fil.Filter(data)
		if err != nil {
			// flush
			p.dataReadPosition = 0
			p.dataWritePosition = 0
			break
		}
		if filterResult == nil {
			// filter 不成功，但没有err，则直接跳出
			break
		}
		p.readPositionAdd(filterResult.GetPackageLength())

		// 事件
		//filterResult.GetPackageBuffer()

		if p.onDataCallback == nil {
			log.Println("no data callback")
			return
		}

		//
		p.onDataCallback(filterResult)
	}
}

// 将读的指针循环右移
func (p *Packet) readPositionAdd(length int) {
	bufLen := p.bufferZoneLength()
	length = length % bufLen
	if span := bufLen - p.dataReadPosition; span <= length {
		firstPartLen := span
		p.dataReadPosition = length - firstPartLen
	} else {
		p.dataReadPosition += length
	}
}

// 将写的指针循环右移
func (p *Packet) writePositionAdd(length int) {
	if length == 0 {
		return
	}

	bufLen := p.bufferZoneLength()
	length = length % bufLen
	if span := bufLen - p.dataWritePosition; span <= length {
		firstPartLen := span
		p.dataWritePosition = length - firstPartLen
	} else {
		p.dataWritePosition += length
	}
}

// 这里有一个临界的问题，主要是 如何表示0值和如何表示满值（即刚好和缓冲区大小一致）
// 目前，用 w - r = 0 则为 0值，即 []byte{}
// 用 w 和 r 相邻（循环相邻 比如 1,2    5,6   0,len）时为最大值
// 但仔细看看，是不是少了一位，比如bufferZone len = 10,那么 r=0,w=9时， dataLength = 9 - 0 =9 ...

// 当前缓冲区内数据的长度
func (p *Packet) currentDataLength() int {
	span := p.dataWritePosition - p.dataReadPosition
	if span >= 0 {
		return span
	}

	return p.bufferZoneLength() - -span
}

// 向缓冲区插入数据
//
// 只管插入数据，不管是否溢出吗？考虑一下
func (p *Packet) insertBuffer(buf []byte) error {

	bufLen := len(buf)
	if bufLen > p.bufferZoneLength() {
		// 这里应该清空数据
		p.dataWritePosition = 0
		p.dataReadPosition = 0
		return errors.New("插入的数据不能大于缓冲区的长度")
	}

	zoneCap := len(p.bufferZone)
	//
	if lengthSpan := p.dataWritePosition + bufLen - zoneCap; lengthSpan >= 0 {
		// 需要截成部分

		// 前半部分，放尾巴
		err := pkgBuffer.BlockCopy(buf, 0, p.bufferZone, p.dataWritePosition, bufLen-lengthSpan)
		if err != nil {
			return err
		}

		// 后半部分，放开头
		err = pkgBuffer.BlockCopy(buf, bufLen-lengthSpan, p.bufferZone, 0, lengthSpan)
		if err != nil {
			return err
		}
		p.dataWritePosition = lengthSpan // 截后的长度
	} else {
		// 不用截
		err := pkgBuffer.BlockCopy(buf, 0, p.bufferZone, p.dataWritePosition, bufLen)
		if err != nil {
			return err
		}
		p.dataWritePosition += bufLen
	}

	return nil
}

// 取当前的数据，即 dataReadPosition 到 dataWritePosition之间的数据
func (p *Packet) getCurrentData() ([]byte, error) {
	length := p.currentDataLength()
	buffer := make([]byte, length, length)
	bufferLength := p.bufferZoneLength()

	// 判断读的部分，是否溢出，如果溢出，取到末尾后，再取开始的部分
	if p.dataReadPosition+length > bufferLength {
		//第一部分，取从readPosition 开始到末尾

		firstLength := bufferLength - p.dataReadPosition
		err := pkgBuffer.BlockCopy(p.bufferZone, p.dataReadPosition, buffer, 0, firstLength)
		if err != nil {
			return []byte{}, err
		}

		// 第二部分，从开始位置开始取，直到取满长度
		err = pkgBuffer.BlockCopy(p.bufferZone, 0, buffer, firstLength, length-firstLength)
		if err != nil {
			return []byte{}, err
		}

	} else {
		err := pkgBuffer.BlockCopy(p.bufferZone, p.dataReadPosition, buffer, 0, length)
		if err != nil {
			return []byte{}, err
		}
	}

	return buffer, nil
}

func (p *Packet) Flush() {
	p.m.Lock()
	p.dataWritePosition = p.dataReadPosition
	p.m.Unlock()

}
