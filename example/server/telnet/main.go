package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/terminator"
	"github.com/lopygo/lopy_socket/tcp/server"
)

func main() {
	c := server.ConfigDefault()

	c.Host = "0.0.0.0"
	c.Port = 18222
	c.PacketFilter, _ = terminator.NewFilter([]byte{0x0a})

	s, err := server.NewServer(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	s.SetOnStarted(func(sender *server.Server, addr *net.TCPAddr) {
		fmt.Println("server started: ", addr.String())
	})
	s.SetOnStopped(func(sender *server.Server) {
		fmt.Println("server stoped: ", sender.GetAddr())
	})

	s.SetOnSessionClose(func(sender *server.Server, session *server.Session) {
		fmt.Println("session close: ", session.SessionID())
	})

	s.SetOnConnected(func(sender *server.Server, session *server.Session) {
		fmt.Println("session connected: ", session.SessionID())
	})

	s.SetOnError(func(sender *server.Server, err error) {
		fmt.Println("show error: ", err)
	})

	s.SetOnSessionError(func(sender *server.Server, session *server.Session, err error) {
		fmt.Println("show session error: [", session.SessionID(), "] ", err)
	})

	s.SetOnData(func(sender *server.Server, session *server.Session, buf []byte) {
		fmt.Printf("received:  % 02X \n", buf)
	})
	s.SetOnDataFiltered(func(sender *server.Server, session *server.Session, dataResult filter.IFilterResult) {
		fmt.Printf("received a package: % 02X \n", dataResult.GetDataBuffer())
	})

	err = s.Start()
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
			s.Stop()
			tick.Stop()
			time.Sleep(time.Second)
			return
		case <-tick.C:

			fmt.Println("show session id: ", len(s.Sessions()))
			for _, v := range s.Sessions() {
				fmt.Println("  - ", v.SessionID())
			}
		}
	}
}
