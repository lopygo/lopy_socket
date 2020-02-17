package packet

import "errors"

// 配置
//
type Option struct {
	// 缓冲区长度
	Length int

	// 数据最大长度，超过这个长度则报异常
	DataMaxLength int

	// 暂时不用
	//// 缓冲区最大长度，（如果启用缓冲区自动扩容，那么这个是指最大扩容后的最大长度）
	//bufferZoneMaxLength int
}

// 检查配置是否正确
// 只做一些基本检查
func (p *Option) Check() error {
	if p.Length < p.DataMaxLength {
		return errors.New("zone length can not lt data length")
	}

	return nil
}

// 默认的option
func NewOptionDefault() *Option {
	return &Option{Length: 4096, DataMaxLength: 512}
}

// 类似于构造函数吧
func NewOption(length int, dataMaxLength int) (*Option, error) {
	option := Option{Length: length, DataMaxLength: dataMaxLength}
	err := option.Check()

	if err != nil {
		return nil, err
	}

	return &option, nil
}
