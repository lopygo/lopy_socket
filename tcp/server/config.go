package server

import "github.com/lopygo/lopy_socket/packet/filter"

type Config struct {
	Host string

	Port uint16

	BufferZoneLenth int

	DataMaxLength int

	Heartbeat int

	PacketFilter filter.IFilter
}

func ConfigDefault() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            12345,
		BufferZoneLenth: 1024,
		DataMaxLength:   1024,
		PacketFilter:    nil,
		Heartbeat:       15,
	}
}
