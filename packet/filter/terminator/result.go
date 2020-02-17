package terminator

import (
	"errors"
	buffer2 "github.com/lopygo/lopy_socket/packet/buffer"
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
	endBuffer := p.creator.GetTerminatorBuffer()
	endBufferLen := len(endBuffer)
	endBufferStart := length - endBufferLen
	if endBufferStart < 0 {
		return errors.New("buffer is shorter than endBuffer")
	}
	tmpBuffer := buffer[endBufferStart:length]
	for k,v := range tmpBuffer {
		if endBuffer[k] != v {
			return errors.New("endBuffer error")
		}
	}

	// 验证完成
	p.SetPackageBuffer(buffer)

	dataBuffer := make([]byte,endBufferStart)
	err := buffer2.BlockCopy(buffer,0,dataBuffer,0,len(dataBuffer))
	if err != nil {
		return err
	}
	p.SetDataBuffer(dataBuffer)

	return nil
}
