package fixed_head

import (
	"bytes"
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
	bufferLength := len(buffer)

	if p.creator == nil {
		return errors.New("creator is not empty")
	}

	err := p.creator.checkBuffer(buffer)
	if err != nil {
		return err
	}

	// bodyLength
	endOffset, err := p.creator.getEndOffset(buffer)
	if err != nil {
		return err
	}

	// 包的总长度不对
	if bufferLength != endOffset {
		return errors.New("length of buffer error")
	}

	// 验证完成
	buf := bytes.Buffer{}
	buf.Write(buffer)
	p.SetPackageBuffer(buf.Bytes())

	buf2 := bytes.Buffer{}
	buf2.Write(buffer[p.creator.bodyOffset:])
	p.SetDataBuffer(buf2.Bytes())

	return nil
}
