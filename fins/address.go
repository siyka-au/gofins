package fins

import "net"

// A FINS device address
type finsAddress struct {
	network byte
	node    byte
	unit    byte
}

// Address A full device address
type Address struct {
	udpAddress  *net.UDPAddr
	finsAddress finsAddress
}

// NewAddress Generates a new FINS address
func NewAddress(ip net.IP, port int, network, node, unit byte) Address {
	return Address{
		udpAddress: &net.UDPAddr{
			IP:   ip,
			Port: port,
		},
		finsAddress: finsAddress{
			network: network,
			node:    node,
			unit:    unit,
		},
	}
}
