package filter

import (
	"errors"
	buffer2 "lopy_socket/packet/buffer"

	//"lopy_socket/packet"

)

//
type TelnetFilterResult struct {
	packageBuffer []byte
	dataBuffer []byte

}

func (p *TelnetFilterResult) PackageLength() int {
	return len(p.packageBuffer)
}

func (p *TelnetFilterResult) PackageBuffer() []byte {
	return p.packageBuffer
}

func (p *TelnetFilterResult) DataBuffer() []byte {
	return p.dataBuffer
}

func (p *TelnetFilterResult) DataLength() int {
	return len(p.dataBuffer)
}

func (p *TelnetFilterResult) SetPackageBuffer(buffer []byte) error {
	// 需要验证一下不
	length := len(buffer)
	if 0x0d != buffer[length-2 ] || 0x0a != buffer[length-1 ] {
		return errors.New("telnet分隔符不对")
	}

	// 验证完成

	p.packageBuffer = buffer

	dataBuffer := make([]byte,p.PackageLength() - 2)
	err := buffer2.BlockCopy(buffer,0,dataBuffer,0,len(dataBuffer))
	if err != nil {
		return err
	}
	p.dataBuffer = dataBuffer
	return nil
}

// telnet
type TelnetFilter struct {

}

func (p *TelnetFilter) GetFilterResult() IFilterResult {
	var filterResult IFilterResult = &TelnetFilterResult{}

	return filterResult

}

func (p TelnetFilter) Filter(buffer []byte) (IFilterResult,error) {
	for k,v := range buffer {
		if k == 0 {
			continue
		}

		// 0x0a 换行 0xod回车
		if v == 0x0a && buffer[k - 1] == 0x0d{
			result := p.GetFilterResult()
			packageBuffer := make([]byte,k+1)
			err := buffer2.BlockCopy(buffer,0, packageBuffer,0,len(packageBuffer))
			if err != nil {
				return nil,err
			}
			err2 := result.SetPackageBuffer(packageBuffer)
			if err2 != nil {
				return nil,err
			}
			return result,nil
		}
	}

	return nil,nil
}
