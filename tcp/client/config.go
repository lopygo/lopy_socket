package client

import "github.com/lopygo/lopy_socket/packet/filter"

type Config struct {
	Ip              string
	Port            uint16
	BufferZoneLenth int

	DataMaxLength int
	// 单位 毫秒
	KeepAlive int
	// 单位 毫秒
	Timeout int
	// 单位 秒
	Heartbeat int

	DataFilter filter.IFilter
}

func ConfigDefault() Config {
	return Config{
		Ip:              "127.0.0.1",
		Port:            12345,
		Timeout:         1000,
		KeepAlive:       2000,
		BufferZoneLenth: 1024,
		DataMaxLength:   512,
		Heartbeat:       15,
	}
}
