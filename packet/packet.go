package packet

import (
	"errors"
	"fmt"
)

type Packet struct {
	dataFilter *FilterInterface

	// 数据读取点，也就是数据在缓冲区的起点index
	dataReadPosition uint

	// 写数据的点，也就是写之前数据在缓冲区的终点
	//dataWritePosition uint
	dataWritePosition uint

	// 数据当前长度，表示正在处理或处理完成时的长度
	dataCurrentLength uint

	// 这个应该是缓冲区
	bufferZone []byte

	// 缓冲区大小，开始表示
	//bufferLength uint

	// 数据包的最大长度，表示如果一个包超过了这个长度，那么将被丢弃
	dataMaxLength uint

	// 缓冲区的最大空间，如果要做自动扩容，那么，这个配置表示扩容后的最大空间，暂时不做自动扩容
	//bufferZoneMaxLength uint
}

func (p *Packet) bufferLength() uint {
	return uint(len(p.bufferZone))
}
// 设置过滤规则，即 粘/拆包的规则
func (p *Packet) SetFilter(filter *FilterInterface) {
	p.dataFilter = filter
}

// 获取过滤规则，即 粘/拆包的规则
func (p *Packet) GetFilter() (*FilterInterface, error) {
	if nil == p.dataFilter {
		return nil, errors.New("filter must be set")
	}
	return p.dataFilter, nil
}

// 获取可用长度，即总的缓冲区长度减当前数据长度
func (p *Packet) GetAvailableLen()  uint {

	return p.bufferLength() - p.dataCurrentLength
}

// 这个看怎么做，先暂时写到这里
func (p *Packet) OnData(buffer []byte)  {
	fmt.Printf("data length is %d", len(buffer))

}

func (p *Packet) Push(data []byte) error {
	// 为空，则不管
	if nil == data || len(data) == 0 {
		return nil
	}

	dataLength := uint(len(data))

	if dataLength > p.dataMaxLength || dataLength + p.dataCurrentLength > p.dataMaxLength {
		msg := "数据长度限制，请跟据实际情况重新配置 dataMaxLength";
		return errors.New(msg)
	}

	//dataStart := 0

	if p.GetAvailableLen() < dataLength {
		// 缓冲区可用空间，总量不足
		return errors.New("缓冲区可用空间，总量不足，请重新设置缓冲区大小")
	}else if p.dataWritePosition + dataLength > p.bufferLength() {
		// 缓冲区尾部长度不足，需要分两截存储，先填满尾部，再把剩余的从0开始

	}

	return nil
}



// 向缓冲区插入数据
// 只管插入数据，不管是否溢出吗？考虑一下
func (p *Packet) insertBuffer(buf []byte) error {
	zoneCap := uint(cap(p.bufferZone))
	bufLen := uint(len(buf))

	if bufLen > p.bufferLength() {
		// 这里应该清空数据
		p.dataCurrentLength = 0
		p.dataWritePosition = 0
		p.dataReadPosition = 0
		return errors.New("插入的数据不能大于缓冲区的长度")
	}

	if lengthSpan := int(p.dataWritePosition + bufLen - zoneCap); lengthSpan > 0 {
		// 需要截成部分

		// 前半部分，放尾巴
		err := BlockCopy(buf, 0, p.bufferZone, int(p.dataWritePosition), int(bufLen)-lengthSpan)
		if err != nil {
			return err
		}

		// 后半部分，放开头
		err = BlockCopy(buf, int(bufLen)-lengthSpan, p.bufferZone, 0, lengthSpan)
		if err != nil {
			return err
		}
		p.dataWritePosition = uint(lengthSpan) // 截后的长度
	} else {
		// 不用截
		err := BlockCopy(buf, 0, p.bufferZone, int(p.dataWritePosition), int(bufLen))
		if err != nil {
			return err
		}
		p.dataWritePosition += bufLen
	}

	return nil
}


