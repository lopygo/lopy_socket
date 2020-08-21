package main

import (
	"github.com/lopygo/lopy_socket/packet"
	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/terminator"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	var bufList = make([][]byte, 0)

	// 0
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 30
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x20, 0x20, 0x33, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 240
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x20, 0x32, 0x34, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 480
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x20, 0x34, 0x38, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 2220
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x32, 0x32, 0x32, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 1220
	// 2330 粘
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x31, 0x32, 0x32, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x32, 0x33, 0x33, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	// 2220 and parity bit (this can insert the next before start)
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x32, 0x32, 0x32, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D, 0x01,
	})

	// 1000
	bufList = append(bufList, []byte{
		0x02, 0x31, 0x30, 0x20, 0x20, 0x20, 0x31, 0x30, 0x30, 0x30, 0x20, 0x20, 0x20, 0x20, 0x30, 0x30, 0x0D,
	})

	opt, _ := packet.NewOption(1024, 1024)
	opt.Filter, _ = terminator.NewFilter([]byte{0x0D})

	pa := packet.NewPacket(opt)

	pa.OnData(func(dataResult filter.IFilterResult) {
		buf := dataResult.GetDataBuffer()
		bufLen := len(buf)

		// get len
		//log.Println(buf[bufLen-16])
		// check start
		if buf[bufLen-16] != 0x02 {
			return
		}
		str := string(buf[bufLen-13 : bufLen-6])
		weight, err := strconv.ParseInt(strings.Trim(str, " "), 10, 64)

		log.Println(weight, err)

	})

	for {

		for _, v := range bufList {
			err := pa.Put(v)
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(time.Duration(20) * time.Millisecond)
	}

}
