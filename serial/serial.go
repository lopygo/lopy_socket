package serial

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/lopygo/lopy_socket/packet"
	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/tarm/serial"
)

type OnHeartPackageGet func(client *Serial) ([]byte, error)
type OnDataFiltered func(client *Serial, dataResult filter.IFilterResult)
type OnData func(client *Serial, buf []byte)
type OnError func(client *Serial, err error)
type OnClosed func(client *Serial)
type OnOpened func(client *Serial)

type Serial struct {

	//
	OnDataFiltered    OnDataFiltered
	OnData            OnData
	OnError           OnError
	OnClosed          OnClosed
	OnOpened          OnOpened
	OnHeartPackageGet OnHeartPackageGet

	config Config

	serialPort *serial.Port

	lockerOpen sync.Mutex

	ctx      context.Context
	ctxClose func()

	// 心跳到期
	heartbeatExpirtReceived time.Time
	heartbeatExpirtSent     time.Time
}

func NewSerial(config *Config) (*Serial, error) {

	s := new(Serial)
	s.config = *config

	return s, nil
}

func (p *Serial) IsStarted() bool {
	return p.ctx != nil
}

func (p *Serial) createPacket() (*packet.Packet, error) {
	lopyOption, err := packet.NewOption(p.config.BufferZoneLength, p.config.DataMaxLength)
	if err != nil {
		return nil, err
	}
	lopyOption.Filter = p.config.PacketFilter
	lopyPacket := packet.NewPacket(lopyOption)
	lopyPacket.OnData(p.triggerOnDataFiltered)

	return lopyPacket, nil
}

func (p *Serial) Open() error {

	p.lockerOpen.Lock()
	defer p.lockerOpen.Unlock()

	if p.IsStarted() {
		return nil
	}
	lopyPacket, err := p.createPacket()
	if err != nil {
		return err
	}

	p.ctx, p.ctxClose = context.WithCancel(context.Background())

	serialPort, err := serial.OpenPort(p.config.toSerialConf())
	if err != nil {
		p.ctx = nil
		p.ctxClose = nil
		return fmt.Errorf("serial [%s] open error: %s", p.config.Name, err)
	}

	// connected
	p.triggerOpenedCallback()

	go func() {
		<-p.ctx.Done()
		serialPort.Close()
		p.ctx = nil
		p.triggerCloseCallback()
	}()

	p.serialPort = serialPort

	go p.listen(serialPort, lopyPacket)

	return nil
}

func (p *Serial) listen(serialPort *serial.Port, lopyPacket *packet.Packet) {
	defer func() {
		p.Close()
	}()

	p.heartbeatUpdateReceived()

	// check connect status
	if p.config.Heartbeat > 0 {
		//
		loopCheckTicker := time.NewTicker(time.Second)
		defer loopCheckTicker.Stop()
		go p.loopCheckStatus(serialPort, loopCheckTicker)
	}

	for p.IsStarted() {

		//

		buf := make([]byte, p.config.BufferZoneLength)
		theLen, err := serialPort.Read(buf)

		if nil != err {
			if err == io.EOF {
				continue
			}

			p.triggerErrorCallback(err)
			return
		}

		// 读取超时时会发生，但貌似是 ioeof
		if theLen == 0 {
			continue
		}
		// 收到成功，就更新
		p.heartbeatUpdateReceived()

		newBuf := make([]byte, theLen)
		copy(newBuf, buf[:theLen])

		p.triggerOnData(newBuf)

		err = lopyPacket.Put(newBuf)
		if err != nil {
			p.triggerErrorCallback(fmt.Errorf("write buffer error: %+v", err))
			lopyPacket.Flush()
		}

	}

}

func (p *Serial) Close() error {
	if p.ctxClose != nil {
		p.ctxClose()
	}
	return nil
}

func (p *Serial) Write(buf []byte) error {
	if p.serialPort == nil {
		return fmt.Errorf("serial is nil")
	}
	_, err := p.serialPort.Write(buf)
	if err == nil {
		p.heartbeatUpdateSent()
	}

	return err
}

func (p *Serial) heartbeatExpireTime(t int) time.Time {
	return time.Now().Add(time.Duration(p.config.Heartbeat+t) * time.Second)
}

func (p *Serial) heartbeatUpdateSent() {
	p.heartbeatExpirtSent = p.heartbeatExpireTime(0)
}

func (p *Serial) heartbeatUpdateReceived() {
	p.heartbeatExpirtReceived = p.heartbeatExpireTime(3)
}

func (p *Serial) triggerOnDataFiltered(dataResult filter.IFilterResult) {
	//data := dataResult.GetDataBuffer()

	//
	if p.OnDataFiltered == nil {
		return

	}

	go p.OnDataFiltered(p, dataResult)
}

func (p *Serial) triggerOnData(buf []byte) {
	//data := dataResult.GetDataBuffer()

	//
	if p.OnData == nil {
		return

	}

	newBuf := make([]byte, len(buf))
	copy(newBuf, buf)

	go p.OnData(p, newBuf)
}

func (p *Serial) triggerOpenedCallback() {

	// connected
	if p.OnOpened == nil {
		return
	}
	go p.OnOpened(p)
}

func (p *Serial) triggerErrorCallback(err error) {
	if p.OnError == nil {
		return
	}

	go p.OnError(p, err)
}

func (p *Serial) triggerCloseCallback() {
	if p.OnClosed == nil {
		return
	}

	go p.OnClosed(p)
}

func (p *Serial) loopCheckStatus(conn *serial.Port, ticker *time.Ticker) {

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

			err = p.Write(buf)
			if err != nil {
				p.triggerErrorCallback(fmt.Errorf("send heartbeat error: %+v", err))
				continue
			}

		}
	}

}
