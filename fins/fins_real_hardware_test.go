package fins

import (
	"testing"
	"time"
)

func testReadWriteRead(localAddr, plcAddr Address) func(t *testing.T) {
	return func(t *testing.T) {
		c, err := NewClient(localAddr, plcAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		c.SetTimeout(time.Duration(500 * time.Millisecond))
 
		var memAddr uint16 = 1000

		read1, err := c.ReadWords(MemoryAreaDMWord, memAddr, 1)
		if err != nil {
			t.Fatal(err)
		}

		var write1 uint16 = 0x2905
		if (write1 == read1[0]) {
			write1 = 0x2307
		}

		c.WriteWords(MemoryAreaDMWord, memAddr, []uint16{write1})

		read2, err := c.ReadWords(MemoryAreaDMWord, memAddr, 1)
		if err != nil {
			t.Fatal(err)
		}
		if (write1 != read2[0]) {
			t.Errorf("Memory address 0x%04x was read and yielded the value 0x%04x. 0x%04x was then written and reading it back yielded 0x%04x which does not match", memAddr, read1[0], write1, read2[0])
		}
	}
}

func TestFINSWithHardware(t *testing.T) {


	localAddr := NewAddress("0.0.0.0", 9600, 0, 30, 0)
	plcAddr := NewAddress("192.168.250.10", 9600, 0, 10, 0)

	t.Run("ReadWriteRead", testReadWriteRead(localAddr, plcAddr))

}
