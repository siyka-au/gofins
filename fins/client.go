package fins

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/siyka-au/bcd"
)

const defaultResponseTimeout = time.Duration(20 * time.Millisecond)

// Client Omron FINS client
type Client struct {
	conn *net.UDPConn
	resp []chan response
	sync.Mutex
	dst             finsAddress
	src             finsAddress
	sid             byte
	closed          bool
	responseTimeout time.Duration
	byteOrder       binary.ByteOrder
}

// ErrIncompatibleMemoryArea Error when the memory area is incompatible with the data type to be read
var ErrIncompatibleMemoryArea = errors.New("The memory area is incompatible with the data type to be read")

// NewClient creates a new Omron FINS client
func NewClient(localAddr, plcAddr Address) (*Client, error) {
	c := new(Client)
	c.dst = plcAddr.finsAddress
	c.src = localAddr.finsAddress
	c.responseTimeout = defaultResponseTimeout
	c.byteOrder = binary.BigEndian

	conn, err := net.DialUDP("udp", localAddr.udpAddress, plcAddr.udpAddress)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	c.resp = make([]chan response, 256) //storage for all responses, sid is byte - only 256 values
	go c.listenLoop()
	return c, nil
}

// SetTimeout Set response timeout duration
// Default value: 20ms.
// A timeout of zero can be used to block indefinitely.
func (c *Client) SetTimeout(duration time.Duration) {
	c.responseTimeout = duration
}

// SetByteOrder Set byte order for data decoding
// Default value: Big endian
func (c *Client) SetByteOrder(byteOrder binary.ByteOrder) {
	c.byteOrder = byteOrder
}

// Close Closes an Omron FINS connection
func (c *Client) Close() {
	c.closed = true
	c.conn.Close()
}

// MemoryAreaReadBytes Reads a string from the PLC memory area
func (c *Client) MemoryAreaReadBytes(memoryArea MemoryArea, address uint16, readCount uint16) ([]byte, error) {
	if checkIsWordMemoryArea(memoryArea) == false {
		return nil, IncompatibleMemoryAreaError{memoryArea}
	}
	command := memoryAreaReadCommand(memoryAddress{memoryArea, address, 0}, readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	return r.data, nil
}

// MemoryAreaReadWords Reads words from the PLC memory area
func (c *Client) MemoryAreaReadWords(memoryArea MemoryArea, address uint16, readCount uint16) ([]uint16, error) {
	data, e := c.MemoryAreaReadBytes(memoryArea, address, readCount)
	if e != nil {
		return nil, e
	}
	wordData := make([]uint16, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		wordData[i] = c.byteOrder.Uint16(data[i*2 : i*2+2])
	}

	return wordData, nil
}

// MemoryAreaReadString Reads a string from the PLC memory area
func (c *Client) MemoryAreaReadString(memoryArea MemoryArea, address uint16, readCount uint16) (string, error) {
	data, e := c.MemoryAreaReadBytes(memoryArea, address, readCount)
	if e != nil {
		return "", e
	}
	n := bytes.IndexByte(data, 0)
	if n == -1 {
		n = len(data)
	}
	return string(data[:n]), nil
}

// MemoryAreaReadBits Reads bits from the PLC memory area
func (c *Client) MemoryAreaReadBits(memoryArea MemoryArea, address uint16, bitOffset byte, readCount uint16) ([]bool, error) {
	if checkIsBitMemoryArea(memoryArea) == false {
		return nil, IncompatibleMemoryAreaError{memoryArea}
	}
	command := memoryAreaReadCommand(memoryAddress{memoryArea, address, bitOffset}, readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	data := make([]bool, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		data[i] = r.data[i]&0x01 > 0
	}

	return data, nil
}

// ReadClock Reads the PLC clock
func (c *Client) ReadClock() (*time.Time, error) {
	r, e := c.sendCommand(clockReadCommand())
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}
	year, _ := bcd.Decode(r.data[0:1])
	if year < 50 {
		year += 2000
	} else {
		year += 1900
	}
	month, _ := bcd.Decode(r.data[1:2])
	day, _ := bcd.Decode(r.data[2:3])
	hour, _ := bcd.Decode(r.data[3:4])
	minute, _ := bcd.Decode(r.data[4:5])
	second, _ := bcd.Decode(r.data[5:6])

	t := time.Date(
		int(year), time.Month(month), int(day), int(hour), int(minute), int(second),
		0, // nanosecond
		time.Local,
	)
	return &t, nil
}

// WriteClock Reads the PLC clock
func (c *Client) WriteClock(time time.Time) error {
	command := clockWriteCommand(0x19, 0x09, 0x18, 0x01, 0x02, 0x03, 0x02)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return e
	}

	return nil
}

// WriteWords Writes words to the PLC data area
func (c *Client) WriteWords(memoryArea MemoryArea, address uint16, data []uint16) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	l := uint16(len(data))
	bts := make([]byte, 2*l, 2*l)
	for i := 0; i < int(l); i++ {
		binary.LittleEndian.PutUint16(bts[i*2:i*2+2], data[i])
	}
	command := memoryAreaWriteCommand(memoryAddress{memoryArea, address, 0}, l, bts)

	return checkResponse(c.sendCommand(command))
}

// WriteString Writes a string to the PLC data area
func (c *Client) WriteString(memoryArea MemoryArea, address uint16, data string) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	stringBytes := make([]byte, 2*len(data), 2*len(data))
	copy(stringBytes, data)

	command := memoryAreaWriteCommand(memoryAddress{memoryArea, address, 0}, uint16((len(data)+1)/2), stringBytes) //TODO: test on real PLC

	return checkResponse(c.sendCommand(command))
}

// WriteBytes Writes bytes array to the PLC data area
func (c *Client) WriteBytes(memoryArea MemoryArea, address uint16, data []byte) error {
	if checkIsWordMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	command := memoryAreaWriteCommand(memoryAddress{memoryArea, address, 0}, uint16(len(data)), data)
	return checkResponse(c.sendCommand(command))
}

// WriteBits Writes bits to the PLC data area
func (c *Client) WriteBits(memoryArea MemoryArea, address uint16, bitOffset byte, data []bool) error {
	if checkIsBitMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	l := uint16(len(data))
	bts := make([]byte, 0, l)
	var d byte
	for i := 0; i < int(l); i++ {
		if data[i] {
			d = 0x01
		} else {
			d = 0x00
		}
		bts = append(bts, d)
	}
	command := memoryAreaWriteCommand(memoryAddress{memoryArea, address, bitOffset}, l, bts)

	return checkResponse(c.sendCommand(command))
}

// SetBit Sets a bit in the PLC data area
func (c *Client) SetBit(memoryArea MemoryArea, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, true)
}

// ResetBit Resets a bit in the PLC data area
func (c *Client) ResetBit(memoryArea MemoryArea, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, false)
}

// ToggleBit Toggles a bit in the PLC data area
func (c *Client) ToggleBit(memoryArea MemoryArea, address uint16, bitOffset byte) error {
	data, e := c.MemoryAreaReadBits(memoryArea, address, bitOffset, 1)
	if e != nil {
		return e
	}
	return c.bitTwiddle(memoryArea, address, bitOffset, !data[0])
}

func (c *Client) bitTwiddle(memoryArea MemoryArea, address uint16, bitOffset byte, value bool) error {
	if checkIsBitMemoryArea(memoryArea) == false {
		return IncompatibleMemoryAreaError{memoryArea}
	}
	mem := memoryAddress{memoryArea, address, bitOffset}
	var byteValue byte = 0x00
	if value == true {
		byteValue = 0x01
	}
	command := memoryAreaWriteCommand(mem, 1, []byte{byteValue})

	return checkResponse(c.sendCommand(command))
}

func checkResponse(r *response, e error) error {
	if e != nil {
		return e
	}
	if r.endCode != EndCodeNormalCompletion {
		return fmt.Errorf("error reported by destination, end code 0x%x", r.endCode)
	}
	return nil
}

func (c *Client) nextHeader() *header {
	sid := c.incrementSid()
	header := defaultCommandHeader(c.src, c.dst, sid)
	return &header
}

func (c *Client) incrementSid() byte {
	c.Lock() //thread-safe sid incrementation
	c.sid++
	sid := c.sid
	c.Unlock()
	c.resp[sid] = make(chan response) //clearing cell of storage for new response
	return sid
}

func (c *Client) sendCommand(command []byte) (*response, error) {
	header := c.nextHeader()
	bts := encodeHeader(*header)
	bts = append(bts, command...)
	_, err := (*c.conn).Write(bts)
	if err != nil {
		return nil, err
	}

	// if response timeout is zero, block indefinitely
	if c.responseTimeout.Nanoseconds() > 0 {
		select {
		case resp := <-c.resp[header.serviceID]:
			return &resp, nil
		case <-time.After(c.responseTimeout):
			return nil, ResponseTimeoutError{c.responseTimeout}
		}
	} else {
		resp := <-c.resp[header.serviceID]
		return &resp, nil
	}
}

func (c *Client) listenLoop() {
	for {
		buf := make([]byte, 2048)
		n, err := bufio.NewReader(c.conn).Read(buf)
		if err != nil {
			// do not complain when connection is closed by user
			if !c.closed {
				log.Fatal(err)
			}
			break
		}

		if n > 0 {
			ans := decodeResponse(buf[:n])
			c.resp[ans.hdr.serviceID] <- ans
		} else {
			log.Println("cannot read response: ", buf)
		}
	}
}

func checkIsWordMemoryArea(memoryArea MemoryArea) bool {
	if memoryArea == MemoryAreaDMWord ||
		memoryArea == MemoryAreaARWord ||
		memoryArea == MemoryAreaHRWord ||
		memoryArea == MemoryAreaWRWord {
		return true
	}
	return false
}

func checkIsBitMemoryArea(memoryArea MemoryArea) bool {
	if memoryArea == MemoryAreaDMBit ||
		memoryArea == MemoryAreaARBit ||
		memoryArea == MemoryAreaHRBit ||
		memoryArea == MemoryAreaWRBit {
		return true
	}
	return false
}

// @ToDo Asynchronous functions

// WriteDataAsync writes to the PLC data area asynchronously
// func (c *Client) WriteDataAsync(startAddr uint16, data []uint16, callback func(resp response)) error {
// 	sid := c.incrementSid()
// 	cmd := writeDCommand(defaultHeader(c.dst, c.src, sid), startAddr, data)
// 	return c.asyncCommand(sid, cmd, callback)
// }
// func (c *Client) asyncCommand(sid byte, cmd []byte, callback func(resp response)) error {
// 	_, err := c.conn.Write(cmd)
// 	if err != nil {
// 		return err
// 	}
// 	asyncResponse(c.resp[sid], callback)
// 	return nil
// }
//
//if callback == nil {
//	p := responseFrame.Payload()			responseFrame := <-c.resp[header.ServiceID()]
//	response := NewResponse(			p := responseFrame.Payload()
//		p.CommandCode(),			response := NewResponse(
//		binary.BigEndian.Uint16(p.Data()[0:2]),				p.CommandCode(),
//		p.Data()[2:])				binary.BigEndian.Uint16(p.Data()[0:2]),
//	return response, nil				p.Data()[2:])
//		return response, nil
//	}
//
// 	go func(frameChannel chan Frame, callback func(*Response)) {
//		responseFrame := <-frameChannel
//		p := responseFrame.Payload()
//		response := NewResponse(
//			p.CommandCode(),
//			binary.BigEndian.Uint16(p.Data()[0:2]),
//			p.Data()[2:])
//		callback(response)
//	}(c.resp[header.ServiceID()], callback)

// func asyncResponse(ch chan response, callback func(r response)) {
// 	if callback != nil {
// 		go func(ch chan response, callback func(r response)) {
// 			ans := <-ch
// 			callback(ans)
// 		}(ch, callback)
// 	}
// }
