package fins

import (
	"encoding/binary"
	"log"
	"net"
)

// Server Omron FINS server (PLC emulator)
type Server struct {
	addr      Address
	conn      *net.UDPConn
	dmarea    []byte
	bitdmarea []byte
	closed    bool
}

const dmAreaSize = 32768

// NewPLCSimulator Creates a new, simulated PLC
func NewPLCSimulator(plcAddr Address) (*Server, error) {
	s := new(Server)
	s.addr = plcAddr
	s.dmarea = make([]byte, dmAreaSize)
	s.bitdmarea = make([]byte, dmAreaSize)

	conn, err := net.ListenUDP("udp", plcAddr.udpAddress)
	if err != nil {
		return nil, err
	}
	s.conn = conn

	go func() {
		var buf [1024]byte
		for {
			rlen, remote, err := conn.ReadFromUDP(buf[:])
			if rlen > 0 {
				req := decodeRequest(buf[:rlen])
				resp := s.handler(req)

				_, err = conn.WriteToUDP(encodeResponse(resp), &net.UDPAddr{IP: remote.IP, Port: remote.Port})
			}
			if err != nil {
				// do not complain when connection is closed by user
				if !s.closed {
					log.Fatal("Encountered error in server loop: ", err)
				}
				break
			}
		}
	}()

	return s, nil
}

// Works with only DM area, 2 byte integers
func (s *Server) handler(r request) response {
	var endCode uint16
	data := []byte{}
	switch r.commandCode {
	case CommandCodeMemoryAreaRead, CommandCodeMemoryAreaWrite:
		memAddr := decodeMemoryAddress(r.data[:4])
		ic := binary.BigEndian.Uint16(r.data[4:6]) // Item count

		switch memAddr.memoryArea {
		case MemoryAreaDMWord:

			if memAddr.address+ic*2 > dmAreaSize { // Check address boundary
				endCode = EndCodeAddressRangeExceeded
				break
			}

			if r.commandCode == CommandCodeMemoryAreaRead { //Read command
				data = s.dmarea[memAddr.address : memAddr.address+ic*2]
			} else { // Write command
				copy(s.dmarea[memAddr.address:memAddr.address+ic*2], r.data[6:6+ic*2])
			}
			endCode = EndCodeNormalCompletion

		case MemoryAreaDMBit:
			if memAddr.address+ic > dmAreaSize { // Check address boundary
				endCode = EndCodeAddressRangeExceeded
				break
			}
			start := memAddr.address + uint16(memAddr.bitOffset)
			if r.commandCode == CommandCodeMemoryAreaRead { //Read command
				data = s.bitdmarea[start : start+ic]
			} else { // Write command
				copy(s.bitdmarea[start:start+ic], r.data[6:6+ic])
			}
			endCode = EndCodeNormalCompletion

		default:
			log.Printf("Memory area is not supported: 0x%04x\n", memAddr.memoryArea)
			endCode = EndCodeNotSupportedByModelVersion
		}

	default:
		log.Printf("Command code is not supported: 0x%04x\n", r.commandCode)
		endCode = EndCodeNotSupportedByModelVersion
	}
	return response{defaultResponseHeader(r.hdr), r.commandCode, endCode, data}
}

// Close Closes the FINS server
func (s *Server) Close() {
	s.closed = true
	s.conn.Close()
}
