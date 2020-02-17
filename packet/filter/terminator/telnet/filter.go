package telnet

import "github.com/lopygo/lopy_socket/packet/filter/terminator"

func NewFilter() (*terminator.Filter, error) {
	fi, err := terminator.NewFilter([]byte{0x0a})
	if err == nil {
		return fi, nil
	}

	return nil, err
}
