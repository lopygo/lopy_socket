package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/terminator"
	"github.com/lopygo/lopy_socket/serial"
)

func main() {
	c := serial.ConfigDefault("COM10", 9600)
	c.Heartbeat = 10

	c.PacketFilter, _ = terminator.NewFilter([]byte{0x0a})

	s, err := serial.NewSerial(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	s.OnClosed = func(client *serial.Serial) {
		fmt.Println("serial close ......")
	}

	s.OnOpened = func(client *serial.Serial) {
		fmt.Println("serial opened")
	}

	s.OnData = func(client *serial.Serial, buf []byte) {
		fmt.Println("received: ", buf)
	}

	s.OnDataFiltered = func(client *serial.Serial, dataResult filter.IFilterResult) {
		fmt.Printf(`
		################## received package ###########################
		package:  % 02x\n
		data   :  % 02x\n
		
		`, dataResult.GetPackageBuffer(), dataResult.GetDataBuffer())
	}

	s.OnError = func(client *serial.Serial, err error) {
		fmt.Println("show err: ", err)
	}

	s.OnHeartPackageGet = func(client *serial.Serial) ([]byte, error) {
		return []byte{0x01, 02}, nil
	}

	err = s.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	tick := time.NewTicker(10 * time.Second)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			s.Close()
			tick.Stop()
			time.Sleep(time.Second)
			return
		case <-tick.C:
			fmt.Println("test tick")
		}
	}
}
