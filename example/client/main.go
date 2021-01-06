package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lopygo/lopy_socket/client"
	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/fixed_head"
)

// you can use sockit to start a server
//
// url: https://github.com/sinpolib/sokit
//
// data example:
// [000100000004FF020100]
// [000100000004FF020102]
// [000100000004FF020103]

func main() {

	ip := flag.String("ip", "127.0.0.1", "ip of remote server")
	port := flag.Uint("port", 502, "port of remote server")

	flag.Parse()

	lengthType, err := fixed_head.NewLengthType(fixed_head.BufferLength2, fixed_head.OrderTypeBigEndian)
	if err != nil {
		fmt.Println(err)
		return
	}

	clientFilter := fixed_head.NewFilter(4, 6, lengthType)
	clientConf := client.ConfigDefault()
	clientConf.BufferZoneLenth = 1024
	clientConf.DataMaxLength = 512
	clientConf.Heartbeat = 15
	clientConf.Ip = *ip
	clientConf.Port = uint16(*port)
	clientConf.DataFilter = clientFilter

	cli, err := client.NewClient(clientConf)

	if err != nil {
		fmt.Println(err)
		return
	}

	cli.OnClosed = func(client *client.Client) {
		fmt.Println("client closed")

		go func() {
			time.Sleep(10 * time.Second)
			cli.Connect()
		}()
	}

	cli.OnConnected = func(client *client.Client) {
		fmt.Println("client connected")
	}

	cli.OnDataReceived = func(client *client.Client, dataResult filter.IFilterResult) {

		fmt.Println(dataResult.GetDataBuffer())
	}

	cli.OnError = func(client *client.Client, err error) {
		// 错误都暂时不处理
		fmt.Println("show error: ", err)
	}

	cli.OnHeartPackageGet = func(client *client.Client) ([]byte, error) {
		return []byte{
			0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0xFF, 0x02, 0x00, 0xC8, 0x00, byte(2),
		}, nil
	}

	err = cli.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:

			return
		}
	}
}
