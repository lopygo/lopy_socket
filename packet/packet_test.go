package packet

import (
	"bytes"
	"fmt"
	"github.com/lopygo/lopy_socket/packet/filter"
	"github.com/lopygo/lopy_socket/packet/filter/fixed_head"
	"github.com/lopygo/lopy_socket/packet/filter/terminator/telnet"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestPacket_BufferZoneLength(t *testing.T) {
	Convey("缓冲区长度", t, func() {
		pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0}
		So(pa.bufferZoneLength(), ShouldEqual, 12)
	})
}

func TestPacket_SetFilter(t *testing.T) {
	//pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0}
	//pa.SetFilter()
	Convey("set filter", t, func() {
		pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0}
		Convey("default it is nil", func() {
			So(pa.dataFilter, ShouldBeNil)
		})
		Convey("when filter set", func() {
			telnetFilter, _ := telnet.NewFilter()
			pa.SetFilter(telnetFilter)
			So(pa.dataFilter, ShouldNotBeNil)
		})
	})
}

func TestPacket_GetFilter(t *testing.T) {
	Convey("get filter", t, func() {
		pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0}
		Convey("default it is nil", func() {
			fil, err := pa.GetFilter()
			So(err, ShouldBeError)
			So(fil, ShouldBeNil)
		})
		Convey("when filter set", func() {
			telnetFilter, _ := telnet.NewFilter()
			pa.SetFilter(telnetFilter)
			fil, err := pa.GetFilter()
			So(err, ShouldBeNil)
			So(fil, ShouldNotBeNil)
			So(fil, ShouldImplement, (*filter.IFilter)(nil))
		})
	})
}

func TestPacket_GetAvailableLen(t *testing.T) {
	Convey("test getAvailableLen", t, func() {
		pa := Packet{bufferZone: make([]byte, 12)}

		Convey("default it is len", func() {
			So(pa.GetAvailableLen(), ShouldEqual, 11)
		})

		type testItem struct {
			readPosition  int
			writePosition int
			expect        int
		}

		arr := []testItem{
			{
				readPosition:  0,
				writePosition: 0,
				expect:        11,
			},
			{
				readPosition:  0,
				writePosition: 11,
				expect:        0,
			},
			{
				readPosition:  0,
				writePosition: 12, // 原则上是不可能的
				expect:        0,
			},
			{
				readPosition:  0,
				writePosition: 4,
				expect:        7,
			},
			{
				readPosition:  4,
				writePosition: 11,
				expect:        4,
			},
			{
				readPosition:  4,
				writePosition: 6,
				expect:        9,
			},
			{
				readPosition:  12,
				writePosition: 0,
				expect:        11,
			},
		}

		for _, item := range arr {
			builder := bytes.Buffer{}
			builder.WriteString("when (read,write) is (")
			builder.WriteString(strconv.Itoa(item.readPosition))
			builder.WriteString(",")
			builder.WriteString(strconv.Itoa(item.writePosition))
			builder.WriteString("), expect (")
			builder.WriteString(strconv.Itoa(item.expect))
			builder.WriteString(")")

			Convey(builder.String(), func() {
				pa.dataReadPosition = item.readPosition
				pa.dataWritePosition = item.writePosition
				So(pa.GetAvailableLen(), ShouldEqual, item.expect)
			})

		}
	})
}

//
func TestInsertBuffer(t *testing.T) {

	pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0}

	Convey("第一次插入数据", t, func() {
		err := pa.insertBuffer([]byte{1, 2, 3})
		So(err, ShouldBeNil)
		So(pa.bufferZone, ShouldResemble, []byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		So(pa.dataWritePosition, ShouldEqual, 3)
	})

	Convey("第二次插入数据", t, func() {
		err := pa.insertBuffer([]byte{5, 6, 7})
		So(err, ShouldBeNil)
		So(pa.bufferZone, ShouldResemble, []byte{1, 2, 3, 5, 6, 7, 0, 0, 0, 0, 0, 0})
		So(pa.dataWritePosition, ShouldEqual, 6)
	})

	Convey("第三次插入数据", t, func() {
		err := pa.insertBuffer([]byte{10, 20, 30, 40, 50, 60, 70})
		So(err, ShouldBeNil)
		So(pa.bufferZone, ShouldResemble, []byte{70, 2, 3, 5, 6, 7, 10, 20, 30, 40, 50, 60})
		So(pa.dataWritePosition, ShouldEqual, 1)
	})

	Convey("第四次插入数据", t, func() {
		err := pa.insertBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
		So(err, ShouldBeNil)
		So(pa.bufferZone, ShouldResemble, []byte{12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
		So(pa.dataWritePosition, ShouldEqual, 1)
	})

	Convey("第五次插入数据", t, func() {
		err := pa.insertBuffer([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
		So(err, ShouldBeError)
		So(pa.dataWritePosition, ShouldEqual, 0)
		So(pa.currentDataLength(), ShouldEqual, 0)
	})

	Convey("第六次插入数据", t, func() {
		err := pa.insertBuffer([]byte{11, 22, 33})
		So(err, ShouldBeNil)
		So(pa.dataWritePosition, ShouldEqual, 3)
		So(pa.bufferZone[0], ShouldEqual, 11)
		So(pa.bufferZone[1], ShouldEqual, 22)
		So(pa.bufferZone[2], ShouldEqual, 33)
	})
}

// 这个应该有问题，writePosition 应该为开区间，即应该为 0 ~ len
func TestWritePositionAdd(t *testing.T) {
	pa := Packet{bufferZone: make([]byte, 12), dataWritePosition: 0, dataReadPosition: 0}
	type testItem struct {
		length           int
		positionExpected int
	}

	Convey("测试写的位置右移", t, func() {

		arr := []testItem{
			testItem{length: 0, positionExpected: 0},
			testItem{length: 1, positionExpected: 1},
			testItem{length: 0, positionExpected: 1},
			testItem{length: 2, positionExpected: 3},
			testItem{length: 3, positionExpected: 6},
			testItem{length: 5, positionExpected: 11}, // 这个
			testItem{length: 1, positionExpected: 0},  // 这个 应该是12，再想一想,应该没有错
			testItem{length: 1, positionExpected: 1},
			testItem{length: 12, positionExpected: 1},
			testItem{length: 13, positionExpected: 2},
			testItem{length: 26, positionExpected: 4},
			//testItem{length: 1, positionExpected: 4}, // 专门写个错的在这里，后续直接删
		}

		for k, item := range arr {
			Convey("number "+strconv.Itoa(k+1)+": move left "+strconv.Itoa(item.length), func() {
				pa.writePositionAdd(item.length)
				So(pa.dataWritePosition, ShouldEqual, item.positionExpected)
			})
		}

	})

}

func TestReadPositionAdd(t *testing.T) {
	pa := Packet{bufferZone: make([]byte, 15), dataWritePosition: 0, dataReadPosition: 0}
	type testItem struct {
		length           int
		positionExpected int
	}

	Convey("测试读的位置右移", t, func() {

		arr := []testItem{
			testItem{length: 0, positionExpected: 0},
			testItem{length: 1, positionExpected: 1},
			testItem{length: 0, positionExpected: 1},
			testItem{length: 2, positionExpected: 3},
			testItem{length: 3, positionExpected: 6},
			testItem{length: 8, positionExpected: 14},
			testItem{length: 1, positionExpected: 0},
			testItem{length: 1, positionExpected: 1},
			testItem{length: 15, positionExpected: 1},
			testItem{length: 16, positionExpected: 2},
			testItem{length: 32, positionExpected: 4},
		}

		for k, item := range arr {
			Convey("number "+strconv.Itoa(k+1)+": move left "+strconv.Itoa(item.length), func() {
				pa.readPositionAdd(item.length)
				So(pa.dataReadPosition, ShouldEqual, item.positionExpected)
			})
		}
	})
}

// 这里应该也有错的
func TestGetCurrentData(t *testing.T) {

	Convey("get current data for (readPosition,writePosition)", t, func() {
		pa := Packet{bufferZone: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, dataWritePosition: 0, dataReadPosition: 0}
		type testItem struct {
			readPosition  int
			writePosition int
			expectBuffer  []byte
			isErr         bool
		}
		arr := []testItem{
			testItem{readPosition: 0, writePosition: 0, expectBuffer: []byte{}},
			testItem{readPosition: 0, writePosition: 1, expectBuffer: []byte{1}},
			testItem{readPosition: 0, writePosition: 2, expectBuffer: []byte{1, 2}},
			testItem{readPosition: 1, writePosition: 0, expectBuffer: []byte{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			testItem{readPosition: 1, writePosition: 1, expectBuffer: []byte{}},
			testItem{readPosition: 1, writePosition: 5, expectBuffer: []byte{2, 3, 4, 5}},
			testItem{readPosition: 3, writePosition: 1, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1}},
			testItem{readPosition: 3, writePosition: 14, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
			//testItem{readPosition: 3, writePosition: 15, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,15}},// 这个和下面这个
			testItem{readPosition: 3, writePosition: 0, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}}, // 这个和上面这个
			testItem{readPosition: 3, writePosition: 1, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1}},
			testItem{readPosition: 3, writePosition: 2, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2}},
			//testItem{readPosition: 3, writePosition: 3, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,15,1,2,3}}, // 这个缺了的
			testItem{readPosition: 3, writePosition: 3, expectBuffer: []byte{}},
			testItem{readPosition: 3, writePosition: 4, expectBuffer: []byte{4}},
			//testItem{readPosition: 3, writePosition: 4, expectBuffer: []byte{41}}, // 故意写错一个先

			// 后面还要加上，尾巴的临界值
			//testItem{readPosition: 1, writePosition: 0, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2}},
			testItem{readPosition: 12, writePosition: 13, expectBuffer: []byte{13}},
			testItem{readPosition: 12, writePosition: 14, expectBuffer: []byte{13, 14}},
			testItem{readPosition: 12, writePosition: 0, expectBuffer: []byte{13, 14, 15}},
			testItem{readPosition: 13, writePosition: 13, expectBuffer: []byte{}},
			testItem{readPosition: 13, writePosition: 14, expectBuffer: []byte{14}},
			testItem{readPosition: 13, writePosition: 0, expectBuffer: []byte{14, 15}},
		}

		for k, item := range arr {
			Convey(fmt.Sprintf("number %d: position (%d,%d)", k+1, item.readPosition, item.writePosition), func() {
				pa.dataReadPosition = item.readPosition
				pa.dataWritePosition = item.writePosition
				data, err := pa.getCurrentData()

				if item.isErr {
					So(err, ShouldBeError)
				} else {
					So(data, ShouldResemble, item.expectBuffer)
				}
			})
		}
	})

}

// 这里应该也有错的
func TestGetCurrentData2(t *testing.T) {
	pa := Packet{bufferZone: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, dataWritePosition: 0, dataReadPosition: 0}
	type testItem struct {
		readPositionAdd  int
		writePositionAdd int
		expectBuffer     []byte
		isErr            bool
	}

	arr := []testItem{
		testItem{readPositionAdd: 0, writePositionAdd: 0, expectBuffer: []byte{}, isErr: false},
		testItem{readPositionAdd: 0, writePositionAdd: 1, expectBuffer: []byte{1}, isErr: false},
		testItem{readPositionAdd: 1, writePositionAdd: 1, expectBuffer: []byte{2}, isErr: false},
		testItem{readPositionAdd: 1, writePositionAdd: 0, expectBuffer: []byte{}, isErr: false},
		testItem{readPositionAdd: 1, writePositionAdd: 0, expectBuffer: []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2}, isErr: false},
	}

	Convey("当前数据开始测试", t, func() {

		for _, item := range arr {
			pa.readPositionAdd(item.readPositionAdd)
			pa.writePositionAdd(item.writePositionAdd)
			data, err := pa.getCurrentData()

			if item.isErr {
				So(err, ShouldBeError)
			} else {
				So(data, ShouldResemble, item.expectBuffer)
			}
		}
	})

}

func TestCurrentDataLength(t *testing.T) {
	pa := Packet{bufferZone: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, dataWritePosition: 0, dataReadPosition: 0}
	type testItem struct {
		readPositionAdd  int
		writePositionAdd int
		expectLength     int
	}

	arr := []testItem{
		{
			readPositionAdd:  0,
			writePositionAdd: 0,
			expectLength:     0,
		},

		{
			readPositionAdd:  0,
			writePositionAdd: 1,
			expectLength:     1,
		},
		{
			readPositionAdd:  1,
			writePositionAdd: 0,
			expectLength:     0,
		},
		{
			readPositionAdd:  0,
			writePositionAdd: 5,
			expectLength:     5,
		},
		{
			readPositionAdd:  0,
			writePositionAdd: 8,
			expectLength:     13,
		},
		{
			// 这里指又超了一圈了
			readPositionAdd:  0,
			writePositionAdd: 5,
			expectLength:     3,
		},
		{
			// read 移4位，正常情况不会遇到，但这是测试，就试试
			readPositionAdd:  4,
			writePositionAdd: 0,
			expectLength:     14,
		},
		{
			// read 移4位，正常情况不会遇到，但这是测试，就试试
			readPositionAdd:  4,
			writePositionAdd: 0,
			expectLength:     10,
		},
	}

	Convey("当前数据长度开始测试", t, func() {

		for _, item := range arr {
			pa.readPositionAdd(item.readPositionAdd)
			pa.writePositionAdd(item.writePositionAdd)
			So(pa.currentDataLength(), ShouldEqual, item.expectLength)
		}
	})
}

func TestNewPacket(t *testing.T) {
	Convey("NewPacket", t, func() {

		thePacket := NewPacket(NewOptionDefault())
		//theFilter,_ := telnet.NewFilter()

		Convey("hehe", func() {
			So(thePacket, ShouldNotBeNil)
		})

	})

}
func TestPacket_Put(t *testing.T) {
	Convey("packet put test", t, func() {
		//telnet := filter.TelnetFilter{}
		telnetFilter, _ := telnet.NewFilter()
		callbackTimes := 0
		pa := Packet{bufferZone: make([]byte, 16), dataWritePosition: 0, dataReadPosition: 0, dataFilter: telnetFilter, dataMaxLength: 1024}

		Convey("when no error", func() {
			buf := []byte{1, 2, 3, 4, 0x0d, 0x0a, 11, 22, 0x0d, 0x0a}
			pa.OnData(func(dataResult filter.IFilterResult) {
				callbackTimes++
			})
			err := pa.Put(buf)
			So(err, ShouldBeNil)
			//done <- true

			Convey("filtered data count should be 2", func() {
				So(callbackTimes, ShouldEqual, 2)
			})
		})

	})

	Convey("put with fixed_head", t, func() {

		type exampleItem struct {
			put    [][]byte
			expect []byte
		}

		exampleList := []exampleItem{
			{
				// 完整的包
				put:    [][]byte{{0, 0, 0, 0, 0, 4, 0, 4, 5, 6}},
				expect: []byte{0, 4, 5, 6},
			},
			{
				//粘包的情况
				put:    [][]byte{{0, 0, 0, 0, 0, 2, 1, 7, 0, 0, 0, 0, 0, 3, 2, 5, 6}},
				expect: []byte{1, 7},
			},
			{
				//粘包，中第二个包的结果，传一个空的串进去
				put:    [][]byte{{}},
				expect: []byte{2, 5, 6},
			},
			{
				//拆包的情况1,
				put:    [][]byte{{0}, {0}, {0}, {0}, {0}, {5}, {3}, {3}, {3}, {4}, {5}},
				expect: []byte{3, 3, 3, 4, 5},
			},
			{
				//拆包的情况2,
				put:    [][]byte{{0, 0}, {0, 0}, {0, 3}, {4, 5}, {9, 1}, {1, 0, 0, 0}, {4, 5}, {0, 0}, {3, 0}, {4}, {5}, {1}, {1}},
				expect: []byte{4, 5, 9},
			},
			{
				//拆包的情况2,的第二个结果
				put:    [][]byte{{}},
				expect: []byte{5, 0, 0, 3},
			},
		}
		// 个数也要对得上
		packageCount := 0
		packageCountExpect := len(exampleList)
		lengthType, _ := fixed_head.NewLengthType(fixed_head.BufferLength4, fixed_head.OrderTypeBigEndian)
		theFilter := fixed_head.NewFilter(2, 6, lengthType)

		thePacket := NewPacket(NewOptionDefault())
		thePacket.SetFilter(theFilter)
		thePacket.OnData(func(dataResult filter.IFilterResult) {
			packageCount++
			dataBuf := dataResult.GetDataBuffer()
			index := int(dataBuf[0])
			Convey("callback "+strconv.Itoa(index), func() {
				So(dataBuf, ShouldResemble, exampleList[index].expect)
			})

		})

		for _, item := range exampleList {
			for _, put := range item.put {

				err := thePacket.Put(put)
				So(err, ShouldBeNil)
			}

		}
		So(packageCount, ShouldEqual, packageCountExpect)
	})
}
