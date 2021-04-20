package fins

import (
	"bufio"
	"log"
	"net"
)

// UDPServerProvider implements ServerProvider interface.
type UDPServerProvider struct {
	conn *net.UDPConn
	resp []chan Frame
	quit chan bool
}

var _ ServerProvider = (*UDPServerProvider)(nil)

func NewUDPServerProvider(plcAddr string) (*UDPServerProvider, error) {
	addr, err := net.ResolveUDPAddr("", plcAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	s := new(UDPServerProvider)
	s.conn = conn
	s.resp = make([]chan Frame, 256) //storage for all responses, sid is byte - only 256 values
	s.quit = make(chan bool)
	go s.listenLoop()
	return s, nil
}

// CloseConnection Closes an Omron FINS connection
func (s *UDPServerProvider) close() error {
	s.quit <- true
	s.conn.Close()
	return nil
}

func (s *UDPServerProvider) listenLoop() {
	for {
		select {
		case <-s.quit:
			return
		default:
			buf := make([]byte, 2048)
			n, err := bufio.NewReader(s.conn).Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			if n > 0 {
				ans := decodeFrame(buf[0:n])
				if err != nil {
					log.Println("failed to parse response: ", err, " \nresponse: ", buf[0:n])
				} else {
					//c.resp[ans.Header.sid] <- *ans
					log.Println("Received: ", ans.Header.sid, " - ", ans.Payload.CommandCode, " - ", ans.Payload.Data)
				}
			} else {
				log.Println("cannot read response: ", buf)
			}
		}
	}
}
