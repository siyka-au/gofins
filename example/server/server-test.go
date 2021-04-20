package main

import (
	"fmt"
	"log"

	fins "github.com/siyka-au/gofins"
)

func main() {

	bindAddr := "192.168.250.3:9600"
	provider, err := fins.NewUDPServerProvider(bindAddr)

	if err != nil {
		log.Fatal()
		panic(fmt.Sprintf("failed to connect to PLC at %v", bindAddr))
	}

	s := fins.NewServer(provider, fins.Address{Network: 0, Node: 3, Unit: 0})
	defer s.Close()

	for {
	}
}
