package client

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

type OnHeartPackageGet func(client *Client) ([]byte, error)

type OnDataReceived func(client *Client, dataResult filter.IFilterResult)

type OnError func(client *Client, err error)
type OnClosed func(client *Client)
type OnConnected func(client *Client)

type Client struct {
	//
	OnDataReceived    OnDataReceived
	OnError           OnError
	OnClosed          OnClosed
	OnConnected       OnConnected
	OnHeartPackageGet OnHeartPackageGet

	// 心跳到期
	heartbeatExpirtReceived time.Time
	heartbeatExpirtSent     time.Time

	config Config
	locker sync.Mutex

	conn net.Conn

	ctx      context.Context
	ctxClose func()
}

func NewClient(conf Config) (*Client, error) {
	cli := new(Client)
	cli.config = conf
	if conf.PacketFilter == nil {
		return nil, fmt.Errorf("client had no filter")
	}
	return cli, nil

}

func (p *Client) IsStarted() bool {
	return nil != p.ctx
}

func (p *Client) Connect() error {

	p.locker.Lock()
	defer p.locker.Unlock()

	if p.IsStarted() {

		return nil
	}

	lopyPacket, err := p.createPacket()
	if err != nil {
		return err
	}

	//
	p.ctx, p.ctxClose = context.WithCancel(context.Background())

	d := net.Dialer{
		Timeout:   time.Duration(p.config.Timeout) * time.Millisecond,
		KeepAlive: time.Duration(p.config.KeepAlive) * time.Millisecond,
	}

	// net.DialTCP()

	conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", p.config.Ip, p.config.Port))
	if err != nil {
		return err
	}

	// connected
	p.triggerConnectedCallback()

	p.conn = conn

	go func() {
		<-p.ctx.Done()
		conn.Close()
		p.ctx = nil
		p.triggerCloseCallback()

	}()

	go p.listen(conn, lopyPacket)

	return nil
}

func (p *Client) createPacket() (*packet.Packet, error) {
	lopyOption, err := packet.NewOption(p.config.BufferZoneLength, p.config.DataMaxLength)
	if err != nil {
		return nil, err
	}
	lopyOption.Filter = p.config.PacketFilter
	lopyPacket := packet.NewPacket(lopyOption)
	lopyPacket.OnData(p.onReceivedData)
	return lopyPacket, nil
}

func (p *Client) listen(conn net.Conn, lopyPacket *packet.Packet) {
	defer func() {
		p.Close()
	}()

	// connected

	p.heartbeatUpdateReceived()

	// check connect status
	if p.config.Heartbeat > 0 {
		//
		loopCheckTicker := time.NewTicker(time.Second)
		defer loopCheckTicker.Stop()
		go p.loopCheckStatus(conn, loopCheckTicker)
	}

	for p.IsStarted() {

		// select {
		// case <-p.ctx.Done():
		// 	// close
		// 	return
		// default:
		// }

		//

		buf := make([]byte, p.config.BufferZoneLength)
		theLen, err := conn.Read(buf)

		if nil != err {
			if err == io.EOF {
				p.triggerErrorCallback(fmt.Errorf("connection closed by remote device"))

			} else {
				p.triggerErrorCallback(err)
			}
			return
		}

		// 这种情况应该不会发生
		// 读取超时时会发生，串口中试过，tcp没有试过
		if theLen == 0 {
			continue
		}
		// 收到成功，就更新
		p.heartbeatUpdateReceived()

		newBuf := make([]byte, theLen)
		copy(newBuf, buf[:theLen])

		err = lopyPacket.Put(newBuf)
		if err != nil {
			p.triggerErrorCallback(fmt.Errorf("write buffer error: %+v", err))
			lopyPacket.Flush()
		}

	}

}

func (p *Client) onReceivedData(dataResult filter.IFilterResult) {
	//data := dataResult.GetDataBuffer()

	//
	if p.OnDataReceived == nil {
		return

	}

	go p.OnDataReceived(p, dataResult)

}

func (p *Client) Close() error {
	if p.ctxClose != nil {
		p.ctxClose()
	}
	return nil
}

func (p *Client) Send(buf []byte) (n int, err error) {
	if p.conn == nil {
		err = fmt.Errorf("client can not connected")
	} else {
		p.heartbeatUpdateSent()
		n, err = p.conn.Write(buf)
	}

	return
}

func (p *Client) loopCheckStatus(conn net.Conn, ticker *time.Ticker) {

	for p.IsStarted() {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:

			//
			// check
			if time.Now().After(p.heartbeatExpirtReceived) {
				p.triggerErrorCallback(fmt.Errorf("heartbeat timeout"))
				p.Close()
				return
			}

			// send for check
			if p.OnHeartPackageGet == nil {
				continue
			}

			if time.Now().Before(p.heartbeatExpirtSent) {
				continue
			}

			buf, err := p.OnHeartPackageGet(p)
			if err != nil {
				p.triggerErrorCallback(fmt.Errorf("get heartbeat template error: %+v", err))
				continue
			}
			if len(buf) == 0 {
				continue
			}

			_, err = p.Send(buf)
			if err != nil {
				p.triggerErrorCallback(fmt.Errorf("send heartbeat error: %+v", err))
				continue
			}

		}
	}

}

func (p *Client) triggerErrorCallback(err error) {
	if p.OnError == nil {
		return
	}

	go p.OnError(p, err)
}

func (p *Client) triggerConnectedCallback() {

	// connected
	if p.OnConnected == nil {
		return
	}
	go p.OnConnected(p)
}

func (p *Client) triggerCloseCallback() {
	if p.OnClosed == nil {
		return
	}

	go p.OnClosed(p)
}

func (p *Client) heartbeatExpireTime(t int) time.Time {
	return time.Now().Add(time.Duration(p.config.Heartbeat+t) * time.Second)
}

func (p *Client) heartbeatUpdateSent() {
	p.heartbeatExpirtSent = p.heartbeatExpireTime(0)
}

func (p *Client) heartbeatUpdateReceived() {
	p.heartbeatExpirtReceived = p.heartbeatExpireTime(3)
}
