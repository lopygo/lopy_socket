package packet

import "errors"

type Option struct {
	// 缓冲区长度
	Length int


	// 数据最大长度，超过这个长度则报异常
	DataMaxLength int


	// 暂时不用
	//// 缓冲区最大长度，（如果启用缓冲区自动扩容，那么这个是指最大扩容后的最大长度）
	//bufferZoneMaxLength int
}

func (option Option) Check() error {
	if option.Length < option.DataMaxLength {
		return errors.New("zone length can not lt data length")
	}

	return nil
}


// 默认的option
func DefaultOption() *Option {
	return &Option{Length:4096,DataMaxLength:512}
}

// 类似于构造函数吧
func GetOption(length int, dataMaxLength int) (*Option,error) {
	option := Option{Length:length,DataMaxLength:dataMaxLength}
	err := option.Check()

	if err != nil {
		return nil,err
	}

	return &option,nil
}
