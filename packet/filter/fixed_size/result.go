package fixed_size

import (
	"errors"
	"github.com/lopygo/lopy_socket/packet/filter"
)

func NewResult(creator *Filter) (*Result, error) {
	res := new(Result)
	// 基本的判断
	if nil == creator {
		return nil, errors.New("creator is not empty")
	}
	res.creator = creator
	return res, nil
}

//
type Result struct {
	filter.Result
	creator *Filter
}

func (p *Result) Assign(buffer []byte) error {
	// 需要验证一下不
	length := len(buffer)

	if p.creator == nil {
		return errors.New("creator is not empty")
	}
	size := p.creator.GetFixedSize()

	if length != size {
		return errors.New("buffer length error")
	}

	// 验证完成
	p.SetPackageBuffer(buffer)
	p.SetDataBuffer(buffer)

	return nil
}
