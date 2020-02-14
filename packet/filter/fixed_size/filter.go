// 固定请求大小的协议
package fixed_size

import (
	"bytes"
	"errors"
	"lopy_socket/packet/filter"
)

type Filter struct {
	// 包的固定大小
	fixedSize int
}


func NewFilter(size int) (*Filter, error) {
	fil := new(Filter)
	// 基本的判断
	if size < 1 {
		return nil, errors.New("size is not empty")
	}
	fil.fixedSize = size
	return fil, nil
}

func (p *Filter) GetFilterResult() (filter.IFilterResult, error) {
	result, err := NewResult(p)

	return result, err

}


func (p *Filter) GetFixedSize() int {
	return p.fixedSize
}


func (p *Filter) Filter(buffer []byte) (filter.IFilterResult, error) {


	// 长度不够直接跳过
	if len(buffer) < p.fixedSize {
		return nil, nil
	}

	result, err := p.GetFilterResult()
	if err != nil {
		return nil, err
	}
	// 有没有其它的方法，比如 copy 之类的
	buf := new(bytes.Buffer)
	buf.Write(buffer[0:p.fixedSize])
	err2 := result.Assign(buf.Bytes())


	if err2 != nil {
		return nil, err2
	}
	return result, nil

}
