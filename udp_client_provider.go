package fins

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
)

// UDPClientProvider implements ClientProvider interface.
type UDPClientProvider struct {
	conn net.Conn
	resp []chan Frame
	quit chan bool
}

var _ ClientProvider = (*UDPClientProvider)(nil)

func NewUDPClientProvider(plcAddr string) (*UDPClientProvider, error) {
	conn, err := net.Dial("udp", plcAddr)

	if err != nil {
		return nil, err
	}

	c := new(UDPClientProvider)
	c.conn = conn
	c.resp = make([]chan Frame, 256) //storage for all responses, sid is byte - only 256 values
	c.quit = make(chan bool)
	go c.listenLoop()
	return c, nil
}

// CloseConnection Closes an Omron FINS connection
func (c *UDPClientProvider) close() error {
	c.quit <- true
	c.conn.Close()
	return nil
}

func (c *UDPClientProvider) sendCommand(header *Header, payload *Payload) (*Response, error) {
	c.resp[header.sid] = make(chan Frame) //clearing cell of storage for new response

	bytes := encodeFrame(NewFrame(header, payload))
	_, err := c.conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	r := <-c.resp[header.sid]
	response := &Response{
		CommandCode: r.Payload.CommandCode,
		EndCode:     binary.BigEndian.Uint16(r.Payload.Data[:2]),
		Data:        r.Payload.Data[2:],
	}
	return response, nil
}

func (c *UDPClientProvider) listenLoop() {
	for {
		select {
		case <-c.quit:
			return
		default:
			buf := make([]byte, 2048)
			n, err := bufio.NewReader(c.conn).Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			if n > 0 {
				ans := decodeFrame(buf[0:n])
				if err != nil {
					log.Println("failed to parse response: ", err, " \nresponse: ", buf[0:n])
				} else {
					c.resp[ans.Header.sid] <- *ans
				}
			} else {
				log.Println("cannot read response: ", buf)
			}
		}
	}
}
