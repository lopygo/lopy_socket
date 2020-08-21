// 固定终止符的协议
package terminator

import (
	"errors"
	"github.com/lopygo/lopy_socket/packet/filter"
)

type Filter struct {
	endBuffer []byte
}

func NewFilter(endBuffer []byte) (*Filter, error) {
	fil := new(Filter)
	// 基本的判断
	if len(endBuffer) == 0 {
		return nil, errors.New("end buffer is not empty")
	}
	fil.endBuffer = endBuffer
	return fil, nil
}

func (p *Filter) GetFilterResult() (filter.IFilterResult, error) {
	result, err := NewResult(p)

	return result, err

}

func (p *Filter) GetTerminatorBuffer() []byte {
	return p.endBuffer
}

// 还可以优化一点点
func (p *Filter) Filter(buffer []byte) (filter.IFilterResult, error) {

	endBufferLen := len(p.endBuffer)
	lastByte := p.endBuffer[endBufferLen-1]

	// 长度不够直接跳过
	if len(buffer) < endBufferLen {
		return nil, nil
	}

	for k, v := range buffer {
		// 如果
		if k < endBufferLen-1 {
			continue
		}

		// 最小的长度了，那么判断最后一位是否相等
		if v != lastByte {
			continue
		}

		// 判断前面的是否相等
		flag := true

		tmpBuffer := buffer[k+1-endBufferLen : k+1]

		for index, dataByte := range tmpBuffer {
			if dataByte != p.endBuffer[index] {
				flag = false
				break
			}
		}

		// 判断flag，没有匹配，继续匹配，直到整个
		if !flag {
			continue
		}

		// 匹配到了

		result, err := p.GetFilterResult()
		if err != nil {
			return nil, err
		}
		packageBuffer := buffer[0 : k+1]
		err2 := result.Assign(packageBuffer)
		if err2 != nil {
			return nil, err2
		}
		return result, nil
	}

	return nil, nil
}
