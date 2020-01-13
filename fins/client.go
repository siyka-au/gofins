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
	"math"

	bcd "github.com/siyka-au/gobcd"
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

// Run Runs program
func (c *Client) Run(opMode CPUOperatingMode) error {
	r, e := c.sendCommand(runCommand(opMode))
	e = checkResponse(r, e)
	if e != nil {
		return e
	}

	return nil
}

// Stop Stops program
func (c *Client) Stop() error {
	r, e := c.sendCommand(stopCommand())
	e = checkResponse(r, e)
	if e != nil {
		return e
	}

	return nil
}

// CPUUnitRead Reads CPU unit information
func (c *Client) CPUUnitRead() error {
	return nil
}

// CPUUnitStatus The status of the CPU unit
type CPUUnitStatus struct {
	Running                   bool
	FlashMemoryWriting        bool
	BatteryPresent            bool
	Standby                   bool
	Mode                      CPUOperatingMode
	FALSError                 bool
	CycleTimeOver             bool
	ProgramError              bool
	IOSettingError            bool
	IOPointOverflow           bool
	FatalInnerBoardError      bool
	DuplicationError          bool
	IOBusError                bool
	MemoryError               bool
	OtherNonFatalError        bool
	SpecialIOUnitSettingError bool
	CS1CPUBusUnitSettingError bool
	BatteryError              bool
	SYSMACBusError            bool
	SpecialIOUnitError        bool
	CPUBusUnitError           bool
	InnerBoardError           bool
	IOVerificationError       bool
	PLCSetupError             bool
	BasicIOUnitError          bool
	InterruptTaskError        bool
	DuplexError               bool
	FALError                  bool
	MessagePresent            [8]bool
	ErrorCode                 uint16
	ErrorMessage              string
}

// CPUUnitStatus CPU unit status
func (c *Client) CPUUnitStatus() (*CPUUnitStatus, error) {
	r, e := c.sendCommand(cpuUnitStatusCommand())
	e = checkResponse(r, e)
	if e != nil {
		return nil, e
	}

	cpuUnitStatus := new(CPUUnitStatus)

	status := r.data[0]
	cpuUnitStatus.Running = status&1 > 0
	cpuUnitStatus.FlashMemoryWriting = status&1<<1 > 0
	cpuUnitStatus.BatteryPresent = status&1<<2 > 0
	cpuUnitStatus.Standby = status&1<<7 > 0

	cpuUnitStatus.Mode = CPUOperatingMode(r.data[1])

	fatalErrorData := binary.BigEndian.Uint16(r.data[2:4])
	cpuUnitStatus.FALSError = fatalErrorData&1<<6 > 0
	cpuUnitStatus.CycleTimeOver = fatalErrorData&1<<8 > 0
	cpuUnitStatus.ProgramError = fatalErrorData&1<<9 > 0
	cpuUnitStatus.IOSettingError = fatalErrorData&1<<10 > 0
	cpuUnitStatus.IOPointOverflow = fatalErrorData&1<<11 > 0
	cpuUnitStatus.FatalInnerBoardError = fatalErrorData&1<<12 > 0
	cpuUnitStatus.DuplicationError = fatalErrorData&1<<13 > 0
	cpuUnitStatus.IOBusError = fatalErrorData&1<<14 > 0
	cpuUnitStatus.MemoryError = fatalErrorData&1<<15 > 0

	nonFatalErrorData := binary.BigEndian.Uint16(r.data[4:6])
	cpuUnitStatus.OtherNonFatalError = nonFatalErrorData&1<<0 > 0
	cpuUnitStatus.SpecialIOUnitSettingError = nonFatalErrorData&1<<2 > 0
	cpuUnitStatus.CS1CPUBusUnitSettingError = nonFatalErrorData&1<<3 > 0
	cpuUnitStatus.BatteryError = nonFatalErrorData&1<<4 > 0
	cpuUnitStatus.SYSMACBusError = nonFatalErrorData&1<<5 > 0
	cpuUnitStatus.SpecialIOUnitError = nonFatalErrorData&1<<6 > 0
	cpuUnitStatus.CPUBusUnitError = nonFatalErrorData&1<<7 > 0
	cpuUnitStatus.InnerBoardError = nonFatalErrorData&1<<8 > 0
	cpuUnitStatus.IOVerificationError = nonFatalErrorData&1<<9 > 0
	cpuUnitStatus.PLCSetupError = nonFatalErrorData&1<<10 > 0
	cpuUnitStatus.BasicIOUnitError = nonFatalErrorData&1<<12 > 0
	cpuUnitStatus.InterruptTaskError = nonFatalErrorData&1<<13 > 0
	cpuUnitStatus.DuplexError = nonFatalErrorData&1<<14 > 0
	cpuUnitStatus.FALError = nonFatalErrorData&1<<15 > 0

	messagePresent := binary.BigEndian.Uint16(r.data[6:8])
	for i := 0; i < len(cpuUnitStatus.MessagePresent); i++ {
		cpuUnitStatus.MessagePresent[i] = messagePresent&1<<i > 0
	}

	cpuUnitStatus.ErrorCode = binary.BigEndian.Uint16(r.data[8:10])

	cpuUnitStatus.ErrorMessage = string(r.data[10:26])

	return cpuUnitStatus, nil
}

// CycleTimeInitialise Initialise cycle time statistics
func (c *Client) CycleTimeInitialise() error {
	r, e := c.sendCommand(cycleTimeInitialiseCommand())
	e = checkResponse(r, e)
	return e
}

// CycleTimeRead Read cycle time statistics
func (c *Client) CycleTimeRead() (*time.Duration, *time.Duration, *time.Duration, error) {
	r, e := c.sendCommand(cycleTimeReadCommand())
	e = checkResponse(r, e)
	if e != nil {
		return nil, nil, nil, e
	}

	avg := time.Duration(binary.BigEndian.Uint32(r.data[0:4]) * 100) * time.Microsecond
	max := time.Duration(binary.BigEndian.Uint32(r.data[4:8]) * 100) * time.Microsecond
	min := time.Duration(binary.BigEndian.Uint32(r.data[8:12]) * 100) * time.Microsecond

	return &avg, &max, &min, nil
}

// ParameterAreaRead Reads data from the PLC parameter area
func (c *Client) ParameterAreaRead(area ParameterArea, address, readCount uint16) (ParameterArea, uint16, uint16, []byte, error) {
	command := parameterAreaReadCommand(area, address, readCount)
	r, e := c.sendCommand(command)
	e = checkResponse(r, e)
	if e != nil {
		return ParameterArea(0), 0, 0, nil, e
	}

	returnedArea := ParameterArea(binary.BigEndian.Uint16(r.data[0:2]))
	returnedBeginningAddress := binary.BigEndian.Uint16(r.data[2:4])
	returnedReadCount := binary.BigEndian.Uint16(r.data[4:6])

	return returnedArea, returnedBeginningAddress, returnedReadCount, r.data[6:], nil
}

// ReadBytes Reads a string from the PLC memory area
func (c *Client) ReadBytes(memoryArea MemoryArea, address uint16, readCount uint16) ([]byte, error) {
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

// ReadWords Reads words from the PLC memory area
func (c *Client) ReadWords(memoryArea MemoryArea, address uint16, readCount uint16) ([]uint16, error) {
	data, e := c.ReadBytes(memoryArea, address, readCount)
	if e != nil {
		return nil, e
	}
	wordData := make([]uint16, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		wordData[i] = c.byteOrder.Uint16(data[i*2 : i*2+2])
	}

	return wordData, nil
}

// ReadString Reads a string from the PLC memory area
func (c *Client) ReadString(memoryArea MemoryArea, address uint16, readCount uint16) (string, error) {
	data, e := c.ReadBytes(memoryArea, address, readCount)
	if e != nil {
		return "", e
	}
	n := bytes.IndexByte(data, 0)
	if n == -1 {
		n = len(data)
	}
	return string(data[:n]), nil
}

// ReadBits Reads bits from the PLC memory area
func (c *Client) ReadBits(memoryArea MemoryArea, address uint16, bitOffset byte, readCount uint16) ([]bool, error) {
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

// WriteClock Sets the PLC clock
func (c *Client) WriteClock(time time.Time) error {

	if time.Year() > 2097 || time.Year() < 1998 {
		return errors.New("Year outside of allowable range of 1998 to 2097 inclusive")
	}

	year := bcd.Encode(uint64(time.Year()))
	month := bcd.Encode(uint64(time.Month()))
	day := bcd.Encode(uint64(time.Day()))
	hour := bcd.Encode(uint64(time.Hour()))
	minute := bcd.Encode(uint64(time.Minute()))
	second := bcd.Encode(uint64(time.Second()))

	command := clockWriteCommand(year[0], month[0], day[0], hour[0], minute[0], second[0], byte(time.Weekday()))
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
		c.byteOrder.PutUint16(bts[i*2:i*2+2], data[i])
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

// WriteFloat32 Writes a float32 as 2 words (4 bytes) to the PLC data area
func (c *Client) WriteFloat32(memoryArea MemoryArea, address uint16, data float32) error {
	buf := make([]byte, 4, 4)
	c.byteOrder.PutUint32(buf[:], math.Float32bits(data))

	return c.WriteBytes(MemoryAreaDMWord, address, buf)
}

// WriteFloat64 Writes a float64 as 4 words (8 bytes) to the PLC data area
func (c *Client) WriteFloat64(memoryArea MemoryArea, address uint16, data float64) error {
	buf := make([]byte, 8, 8)
	c.byteOrder.PutUint64(buf[:], math.Float64bits(data))

	return c.WriteBytes(MemoryAreaDMWord, address, buf)
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
	data, e := c.ReadBits(memoryArea, address, bitOffset, 1)
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

func checkResponse(r *response, e error) error {
	if e != nil {
		return e
	}
	if r.endCode != EndCodeNormalCompletion {
		return fmt.Errorf("error reported by destination, end code 0x%x", r.endCode)
	}
	return nil
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
