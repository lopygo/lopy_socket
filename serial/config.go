package serial

import (
	"time"

	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/tarm/serial"
)

type StopBits byte
type Parity byte

const (
	Stop1     StopBits = 1
	Stop1Half StopBits = 15
	Stop2     StopBits = 2
)

const (
	ParityNone  Parity = 'N'
	ParityOdd   Parity = 'O'
	ParityEven  Parity = 'E'
	ParityMark  Parity = 'M' // parity bit is always 1
	ParitySpace Parity = 'S' // parity bit is always 0
)

type Config struct {
	Name string

	// 码率
	Baud int

	// 数据位 //这个是不是没有用
	DataBits int

	BufferZoneLength int

	DataMaxLength int

	PacketFilter filter.IFilter

	Heartbeat int

	// Size is the number of data bits. If 0, DefaultSize is used.
	Size byte

	// Parity is the bit to use and defaults to ParityNone (no parity bit).
	Parity Parity

	// Number of stop bits to use. Default is 1 (1 stop bit).
	StopBits StopBits

	ReadTimeout int
}

func (p *Config) toSerialConf() *serial.Config {
	i := serial.Config{
		Name: p.Name,
		Baud: p.Baud,
		Size: p.Size,

		// Parity is the bit to use and defaults to ParityNone (no parity bit).
		Parity: serial.Parity(p.Parity),

		// Number of stop bits to use. Default is 1 (1 stop bit).
		StopBits: serial.StopBits(p.StopBits),
	}

	if p.ReadTimeout > 0 {
		i.ReadTimeout = time.Millisecond * time.Duration(p.ReadTimeout)
	}

	return &i
}

func ConfigDefault(name string, baud int) *Config {
	i := new(Config)
	i.Baud = baud
	i.Name = name

	i.BufferZoneLength = 1024
	i.DataMaxLength = 512
	i.Size = 0
	i.Parity = ParityNone
	i.StopBits = Stop1
	i.ReadTimeout = 0
	i.Heartbeat = 0
	return i
}
