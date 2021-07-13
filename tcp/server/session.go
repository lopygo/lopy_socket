package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lopygo/lopy_socket/service"
)

type Session struct {
	connectedTime    time.Time
	lastSentTime     time.Time
	lastReceivedTime time.Time
	sessionID        int64

	data sync.Map

	conn net.Conn

	onClose func(*Session)
}

func newSession(conn net.Conn) (session *Session, err error) {
	session = new(Session)
	session.data = sync.Map{}

	session.connectedTime = time.Now()
	session.lastSentTime = time.Now()
	session.lastReceivedTime = time.Now()

	session.conn = conn
	session.sessionID = service.GetSnowflakeId().Int64()
	return
}

func (receiver *Session) SessionID() int64 {

	return receiver.sessionID
}

func (receiver *Session) Send(buf []byte) (n int, err error) {
	if receiver.conn == nil {
		err = fmt.Errorf("connect can not exist or disconnect")
		return
	}

	receiver.lastSentTime = time.Now()
	n, err = receiver.conn.Write(buf)
	return
}

func (receiver *Session) received(buf []byte) (err error) {

	receiver.lastReceivedTime = time.Now()
	return
}

func (receiver *Session) Close() (err error) {
	if receiver.conn != nil {
		receiver.conn.Close()
	}
	if receiver.onClose != nil {
		receiver.onClose(receiver)
	}
	return
}

func (receiver *Session) SetData(key, value interface{}) {
	receiver.data.Store(key, value)
}

func (receiver *Session) GetData(key interface{}) (interface{}, bool) {

	return receiver.data.Load(key)
}

func (receiver *Session) ConnectedTime() time.Time {

	return receiver.connectedTime
}

func (receiver *Session) LastSentTime() time.Time {

	return receiver.lastSentTime
}

func (receiver *Session) LastReceivedTime() time.Time {

	return receiver.lastReceivedTime
}
