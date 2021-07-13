package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/lopygo/lopy_socket/packet"
	"github.com/lopygo/lopy_socket/packet/filter"
)

const (
	ValueKeyIPAddr CtxValueKey = "ipAddr"
)

type CtxValueKey string

type OnStarted func(sender *Server, addr *net.TCPAddr)
type OnStopped func(sender *Server)
type OnConnected func(sender *Server, session *Session)
type OnSessionClose func(sender *Server, session *Session)
type OnTimeout func(sender *Server, session *Session)
type OnError func(sender *Server, err error)
type OnSessionError func(sender *Server, session *Session, err error)
type OnData func(sender *Server, session *Session, buf []byte)
type OnDataFiltered func(sender *Server, session *Session, dataResult filter.IFilterResult)

type Server struct {
	config Config

	locker sync.Mutex

	connList sync.Map

	ctx      context.Context
	ctxClose func()

	// on event
	onConnected    OnConnected
	onSessionClose OnSessionClose
	onStarted      OnStarted
	onStopped      OnStopped
	onTimeout      OnTimeout
	onError        OnError
	onSessionError OnSessionError
	onData         OnData
	onDataFiltered OnDataFiltered
}

func NewServer(conf Config) (server *Server, err error) {

	server = &Server{
		config: conf,
	}

	server.connList = sync.Map{}

	return
}

func (receiver *Server) SetOnConnected(onConnected OnConnected) {
	receiver.onConnected = onConnected
}
func (receiver *Server) SetOnSessionClose(onClose OnSessionClose) {
	receiver.onSessionClose = onClose
}
func (receiver *Server) SetOnTimeout(onTimeout OnTimeout) {
	receiver.onTimeout = onTimeout
}
func (receiver *Server) SetOnError(onError OnError) {
	receiver.onError = onError
}
func (receiver *Server) SetOnStarted(onStarted OnStarted) {
	receiver.onStarted = onStarted
}
func (receiver *Server) SetOnStopped(onStopped OnStopped) {
	receiver.onStopped = onStopped
}
func (receiver *Server) SetOnSessionError(onSessionError OnSessionError) {
	receiver.onSessionError = onSessionError
}
func (receiver *Server) SetOnData(onData OnData) {
	receiver.onData = onData
}
func (receiver *Server) SetOnDataFiltered(onDataFiltered OnDataFiltered) {
	receiver.onDataFiltered = onDataFiltered
}

func (receiver *Server) GetAddr() *net.TCPAddr {
	// return receiver.
	if receiver.ctx == nil {
		return nil
	}

	v := receiver.ctx.Value(ValueKeyIPAddr)
	if v == nil {
		return nil
	}

	vv, ok := v.(*net.TCPAddr)
	if !ok {
		return nil
	}

	return vv

}

func (receiver *Server) Sessions() []*Session {
	l := make([]*Session, 0)
	receiver.connList.Range(func(key, value interface{}) bool {
		k, ok := value.(*Session)
		if !ok {
			receiver.connList.Delete(key)
			// on error
			// 理论上不会发生
			receiver.triggerOnError(fmt.Errorf("server invalid error"))
		}

		l = append(l, k)
		return true
	})

	return l
}

func (receiver *Server) GetSession(id int64) *Session {
	v, ok := receiver.connList.Load(id)
	if !ok {
		return nil
	}

	k, ok := v.(*Session)
	if !ok {
		receiver.connList.Delete(id)
		// on error
		// 理论上不会发生
		receiver.triggerOnError(fmt.Errorf("get session error: invalid error"))
	}

	return k
}

func (receiver *Server) IsStarted() bool {
	return receiver.ctx != nil
}

func (receiver *Server) Start() (err error) {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	if receiver.IsStarted() {

		return
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", receiver.config.Host, receiver.config.Port))
	if err != nil {
		return
	}

	receiver.ctx, receiver.ctxClose = context.WithCancel(context.Background())
	receiver.ctx = context.WithValue(receiver.ctx, ValueKeyIPAddr, tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		receiver.ctx = nil
		receiver.ctxClose = nil
		return
	}

	go func() {
		<-receiver.ctx.Done()
		listener.Close()
		receiver.triggerOnStoped()
	}()

	// go func() {
	// 	<-receiver.ctx.Done()
	// 	for _, v := range receiver.Sessions() {
	// 		v.Close()
	// 	}
	// }()
	// on start

	receiver.triggerOnStarted(tcpAddr)

	// listen
	go receiver.listen(listener)

	return
}

func (receiver *Server) listen(listener net.Listener) {

	for receiver.IsStarted() {

		// 这段select 貌似没有用，因为accept 阻塞了
		// select {
		// case <-receiver.ctx.Done():
		// 	fmt.Println("server: closing ....")
		// 	return
		// default:
		// 	// default
		// 	// break
		// }

		//
		conn, err := listener.Accept()
		if err != nil {
			// show error ?
			return
		}

		go receiver.connProcess(conn)

	}
}

func (receiver *Server) connProcess(conn net.Conn) {

	session, _ := newSession(conn)
	sessionId := session.SessionID()

	defer func() {
		session.Close()
	}()

	// session list
	receiver.connList.Store(sessionId, session)
	session.onClose = func(s *Session) {
		receiver.triggerOnSessionClose(session)
		receiver.connList.Delete(sessionId)
	}

	receiver.triggerOnConnected(session)

	// init

	// 心跳
	// setTimeout(conn, receiver.config.Heartbeat)

	packetObj := receiver.createPacket()
	packetObj.OnData(func(dataResult filter.IFilterResult) {
		receiver.triggerOnDataFiltered(session, dataResult)
	})

	for receiver.IsStarted() {

		//  这段select 貌似没有用，因为 Read 阻塞了，暂时岛
		// select {
		// case <-receiver.ctx.Done():
		// 	fmt.Println("session closing ... ")
		// 	return
		// default:

		// }

		buf := make([]byte, receiver.config.BufferZoneLenth)
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Println("close from remote")
				return
			}

			if err, ok := err.(net.Error); ok && err.Timeout() {
				receiver.triggerOnTimeout(session)
				return
			}

			receiver.triggerOnSessionError(session, err)
			return
		}

		// tcp 还没有遇到n=0的情况，不知道超时是怎么算的
		if n == 0 {
			time.Sleep(time.Millisecond)
			continue
		}

		bufReal := buf[:n]

		// 心跳
		// setTimeout(conn, p.config.Heartbeat)
		session.received(bufReal)
		receiver.triggerOnData(session, bufReal)

		if err := packetObj.Put(bufReal); err != nil {
			receiver.triggerOnSessionError(session, err)
			packetObj.Flush()
		}

	}
}

func (receiver *Server) createPacket() *packet.Packet {
	op := packet.NewOptionDefault()
	op.Filter = receiver.config.PacketFilter
	op.DataMaxLength = receiver.config.DataMaxLength
	op.Length = receiver.config.BufferZoneLenth

	p := packet.NewPacket(op)

	return p
}

func (receiver *Server) Stop() error {
	if receiver.ctxClose != nil {
		receiver.ctxClose()
	}
	return nil
}

func (receiver *Server) triggerOnConnected(session *Session) {
	if receiver.onConnected == nil {
		return
	}

	go receiver.onConnected(receiver, session)
}

func (receiver *Server) triggerOnSessionClose(session *Session) {
	if receiver.onSessionClose == nil {
		return
	}

	go receiver.onSessionClose(receiver, session)
}

func (receiver *Server) triggerOnTimeout(session *Session) {
	if receiver.onTimeout == nil {
		return
	}

	go receiver.onTimeout(receiver, session)
}

func (receiver *Server) triggerOnError(err error) {
	if receiver.onError == nil {
		return
	}

	go receiver.onError(receiver, err)
}

func (receiver *Server) triggerOnStarted(addr *net.TCPAddr) {
	if receiver.onStarted == nil {
		return
	}

	go receiver.onStarted(receiver, addr)
}

func (receiver *Server) triggerOnStoped() {
	if receiver.onStopped == nil {
		return
	}

	go receiver.onStopped(receiver)
}

func (receiver *Server) triggerOnSessionError(session *Session, err error) {
	if receiver.onSessionError == nil {
		return
	}

	go receiver.onSessionError(receiver, session, err)
}

func (receiver *Server) triggerOnData(session *Session, buf []byte) {
	if receiver.onData == nil {
		return
	}
	newBuf := make([]byte, len(buf))
	copy(newBuf, buf)

	go receiver.onData(receiver, session, newBuf)

}

func (receiver *Server) triggerOnDataFiltered(session *Session, filterResult filter.IFilterResult) {
	if receiver.onDataFiltered == nil {
		return
	}

	go receiver.onDataFiltered(receiver, session, filterResult)
}
