package main

import (
	"fmt"
	"log"
	"time"

	fins "github.com/siyka-au/gofins"
)

func main() {

	plcAddr := "192.168.250.10:9600"
	provider, err := fins.NewUDPClientProvider(plcAddr)

	if err != nil {
		log.Fatal()
		panic(fmt.Sprintf("failed to connect to PLC at %v", plcAddr))
	}

	c := fins.NewClient(provider, fins.Address{
		Network: 0,
		Node:    10,
		Unit:    0,
	}, fins.Address{
		Network: 0,
		Node:    2,
		Unit:    0,
	})
	defer c.Close()

	// z, _ := c.ReadWords(fins.MemoryAreaDMWord, 24000, 2)
	// fmt.Println(z)

	// s, _ := c.ReadString(fins.MemoryAreaDMWord, 10000, 10)
	// fmt.Println(s)
	// fmt.Println(len(s))

	t, _ := c.ReadClock()
	fmt.Println(t.Format(time.RFC3339))

	// b, _ := c.ReadBits(fins.MemoryAreaDMWord, 10473, 2, 1)
	// fmt.Println(b)
	// fmt.Println(len(b))

	// c.WriteWords(fins.MemoryAreaDMWord, 24000, []uint16{z[0] + 1, z[1] - 1})
	// c.WriteBits(fins.MemoryAreaDMBit, 24002, 0, []bool{false, false, false, true,
	// 	true, false, false, true,
	// 	false, false, false, false,
	// 	true, true, true, true})
	// c.SetBit(fins.MemoryAreaDMBit, 24003, 1)
	// c.ResetBit(fins.MemoryAreaDMBit, 24003, 0)
	// c.ToggleBit(fins.MemoryAreaDMBit, 24003, 2)
	c.WriteString(fins.MemoryAreaDMWord, 10000, 10, "WeLoveGoLang!")

	t, _ = c.ReadClock()
	fmt.Println(t.Format(time.RFC3339))
}
