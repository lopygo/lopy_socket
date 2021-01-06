package client

import (
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

	// 最后一次收到data的时间
	lastReceivedTime time.Time

	status bool

	config Config
	locker sync.Mutex

	conn net.Conn
}

func (p *Client) Connect() error {
	if p.status {
		return nil
	}

	go p.listen()

	return nil
}

func (p *Client) triggerErrorCallback(err error) {
	if p.OnError == nil {
		return
	}

	go p.OnError(p, err)
}

func (p *Client) listen() {
	p.locker.Lock()
	defer func() {
		p.Close()
		p.locker.Unlock()
	}()

	if p.status == true {
		return
	}
	p.status = true

	// buffer zone
	lopyOption, err := packet.NewOption(p.config.BufferZoneLenth, p.config.DataMaxLength)
	if err != nil {
		p.triggerErrorCallback(fmt.Errorf("init buffer zone err: %+v", err))
		return
	}
	lopyOption.Filter = p.config.DataFilter
	lopyPacket := packet.NewPacket(lopyOption)
	lopyPacket.OnData(p.onReceivedData)

	// connected
	d := net.Dialer{
		Timeout:   time.Duration(p.config.Timeout) * time.Millisecond,
		KeepAlive: time.Duration(p.config.KeepAlive) * time.Millisecond,
	}

	conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", p.config.Ip, p.config.Port))
	defer func() {
		if p.OnClosed != nil {
			p.OnClosed(p)
		}
	}()
	if err != nil {
		p.triggerErrorCallback(fmt.Errorf("connect failed, err: %+v", err))
		return
	}
	p.conn = conn
	p.lastReceivedTime = time.Now()

	// connected
	if p.OnConnected != nil {
		go p.OnConnected(p)
	}

	// check connect status
	if p.config.Heartbeat > 0 {
		loopCheckTicker := time.NewTicker(time.Duration(p.config.Heartbeat) * time.Second)
		defer loopCheckTicker.Stop()
		go p.loopCheckStatus(conn, loopCheckTicker)
	}

	for p.status {
		buf := make([]byte, 1024)
		theLen, err := conn.Read(buf)

		if nil != err {
			if err == io.EOF {
				p.triggerErrorCallback(fmt.Errorf("connection closed by remote device"))
				break
			}
			// 报错直接关闭
			p.triggerErrorCallback(fmt.Errorf("read buffer error: %+v", err))
			break
		}

		// 这种情况应该不会发生
		if theLen == 0 {
			continue
		}
		// 收到成功，就更新
		p.lastReceivedTime = time.Now()

		newBuf := make([]byte, theLen)
		copy(newBuf, buf)

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
	p.status = false
	if p.conn != nil {

		p.conn.Close()
	}
	return nil
}

func (p *Client) loopCheckStatus(conn net.Conn, ticker *time.Ticker) {

	for p.status {
		select {
		case <-ticker.C:

			//
			timeExpire := time.Now().Add(time.Duration(-2-p.config.Heartbeat) * time.Second)

			if p.lastReceivedTime.Before(timeExpire) {
				p.triggerErrorCallback(fmt.Errorf("heartbeat timeout"))
				p.Close()
				return
			}

			if p.OnHeartPackageGet == nil {
				continue
			}

			buf, err := p.OnHeartPackageGet(p)
			if err != nil {
				p.triggerErrorCallback(fmt.Errorf("get heartbeat template error: %+v", err))
				continue
			}
			_, err = conn.Write(buf)
			if err != nil {
				p.triggerErrorCallback(fmt.Errorf("send heartbeat error: %+v", err))
				continue
			}

		default:

			break
		}
	}

}

func NewClient(conf Config) (*Client, error) {
	cli := new(Client)
	cli.config = conf
	if conf.DataFilter == nil {
		return nil, fmt.Errorf("client had no filter")
	}
	return cli, nil

}
