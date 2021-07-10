package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/terminator"
	"github.com/lopygo/lopy_socket/tcp/client"
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

	clientFilter, err := terminator.NewFilter([]byte{0x0d, 0x0a})
	if err != nil {
		fmt.Println("filter error: ", err)
		return
	}
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

	}

	cli.OnConnected = func(client *client.Client) {
		fmt.Println("client connected")
	}

	cli.OnDataReceived = func(client *client.Client, dataResult filter.IFilterResult) {

		fmt.Println(string(dataResult.GetDataBuffer()))
	}

	cli.OnError = func(client *client.Client, err error) {
		// 错误都暂时不处理
		fmt.Println("show error: ", err)
	}

	cli.OnHeartPackageGet = func(client *client.Client) ([]byte, error) {
		return []byte("AT+STACH0=?\r\n"), nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second * 20)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := cli.Connect()
				if nil == err {
					continue
				}

				fmt.Println("connect error: ", err)
			case <-ctx.Done():

				return
			}
		}
	}()

	go func() {

		err = cli.Connect()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			cancel()
			return
		}
	}
}
