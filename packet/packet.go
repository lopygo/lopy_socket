package packet

import (
	"errors"
	"fmt"
	pkgBuffer "lopy_socket/packet/buffer"
	"lopy_socket/packet/filter"
)

type Packet struct {
	dataFilter filter.IFilter

	// 数据读取点，也就是数据在缓冲区的起点index
	dataReadPosition uint

	// 写数据的点，也就是写之前数据在缓冲区的终点
	dataWritePosition uint

	// 这个应该是缓冲区
	bufferZone []byte

	// 数据包的最大长度，表示如果一个包超过了这个长度，那么将被丢弃
	dataMaxLength uint

	// 缓冲区的最大空间，如果要做自动扩容，那么，这个配置表示扩容后的最大空间，暂时不做自动扩容
	//bufferZoneMaxLength uint
}

// 缓冲区长度
func (p *Packet) bufferZoneLength() uint {
	return uint(len(p.bufferZone))
}

// 设置过滤规则，即 粘/拆包的规则
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
func (p *Packet) GetAvailableLen() uint {
	return p.bufferZoneLength() - p.currentDataLength()
}

// 这个看怎么做，先暂时写到这里
func (p *Packet) OnData(buffer []byte) {
	fmt.Printf("data length is %d", len(buffer))
}

// 外部写入数据，估计要考虑并发问题
func (p *Packet) Put(data []byte) error {
	// 为空，则不管
	if nil == data || len(data) == 0 {
		return nil
	}

	dataLength := uint(len(data))
	if dataLength > p.dataMaxLength {
		msg := "数据长度限制，请跟据实际情况重新配置 dataMaxLength"
		return errors.New(msg)
	}
	if dataLength+p.currentDataLength() > p.bufferZoneLength() {
		msg := "缓冲区大小不足"
		return errors.New(msg)
	}

	// 插入数据
	err := p.insertBuffer(data)
	if err != nil {
		return err
	}

	// filter
	p.readByFilter()

	return nil
}

func (p *Packet) readByFilter() {
	for{
		data,err := p.getCurrentData()
		if err != nil {
			break
		}
		if len(data) == 0 {
			break
		}

		// 从这里开始，的异常，要清空数据
		if p.dataFilter == nil {
			p.dataReadPosition = 0
			p.dataWritePosition = 0
			break
		}
		filterResult,err :=p.dataFilter.Filter(data)
		if err != nil {
			p.dataReadPosition = 0
			p.dataWritePosition = 0
			break
		}
		p.readPositionAdd(filterResult.PackageLength())

		// 事件
		//filterResult.PackageBuffer()

		fmt.Println("??????????????????/")
		fmt.Println(filterResult.PackageBuffer())
		fmt.Println(filterResult.DataBuffer())

	}
}

// 将读的指针循环右移
func (p *Packet) readPositionAdd(length uint) {
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
func (p *Packet) writePositionAdd(length uint) {
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

func (p *Packet) currentDataLength() uint {
	span := int(p.dataWritePosition - p.dataReadPosition)
	if span >= 0 {
		return uint(span)
	}

	return p.bufferZoneLength() - uint(-span)
}

// 向缓冲区插入数据
// 只管插入数据，不管是否溢出吗？考虑一下
func (p *Packet) insertBuffer(buf []byte) error {
	zoneCap := uint(cap(p.bufferZone))
	bufLen := uint(len(buf))

	if bufLen > p.bufferZoneLength() {
		// 这里应该清空数据
		p.dataWritePosition = 0
		p.dataReadPosition = 0
		return errors.New("插入的数据不能大于缓冲区的长度")
	}

	if lengthSpan := int(p.dataWritePosition + bufLen - zoneCap); lengthSpan > 0 {
		// 需要截成部分

		// 前半部分，放尾巴
		err := pkgBuffer.BlockCopy(buf, 0, p.bufferZone, int(p.dataWritePosition), int(bufLen)-lengthSpan)
		if err != nil {
			return err
		}

		// 后半部分，放开头
		err = pkgBuffer.BlockCopy(buf, int(bufLen)-lengthSpan, p.bufferZone, 0, lengthSpan)
		if err != nil {
			return err
		}
		p.dataWritePosition = uint(lengthSpan) // 截后的长度
	} else {
		// 不用截
		err := pkgBuffer.BlockCopy(buf, 0, p.bufferZone, int(p.dataWritePosition), int(bufLen))
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

		firstLength := int(bufferLength - p.dataReadPosition)
		err := pkgBuffer.BlockCopy(p.bufferZone, int(p.dataReadPosition), buffer, 0, firstLength)
		if err != nil {
			return []byte{}, err
		}

		// 第二部分，从开始位置开始取，直到取满长度
		err = pkgBuffer.BlockCopy(p.bufferZone, 0, buffer, firstLength, int(length)-firstLength)
		if err != nil {
			return []byte{}, err
		}

	} else {
		err := pkgBuffer.BlockCopy(p.bufferZone, int(p.dataReadPosition), buffer, 0, int(length))
		if err != nil {
			return []byte{}, err
		}
	}

	return buffer, nil
}
