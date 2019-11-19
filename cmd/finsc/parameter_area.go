package main

import (
	"fmt"
	"encoding/hex"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

type parameterAreaCommand struct {
	fins             **fins.Client
	beginningAddress uint16
	readCount        uint16
	area             string
	file             string
}

func configureParameterArea(app *kingpin.Application, finsc *finscContext) {
	parameterArea := app.Command("parameter-area", "Parameter area")
	pac := &parameterAreaCommand{}
	pac.fins = &finsc.client

	parameterAreaRead := parameterArea.Command("read", "Read parameter area").Action(pac.readParameterArea)
	parameterAreaRead.Arg("area", "Parameter area to read").Default("plc").EnumVar(&pac.area, "plc", "io", "routing", "bus", "unknown")
	parameterAreaRead.Arg("file", "File to save data to").StringVar(&pac.file)
	parameterAreaRead.Flag("beginning-address", "address to start reading data from").Default("0x0000").Uint16Var(&pac.beginningAddress)
	parameterAreaRead.Flag("read-count", "number of words to read").Default("0x0001").Uint16Var(&pac.readCount)

	parameterAreaWrite := parameterArea.Command("write", "Write parameter area").Action(pac.writeParameterArea)
	parameterAreaWrite.Arg("area", "Parameter area to write").Required().EnumVar(&pac.area, "plc", "io", "routing", "bus")
	parameterAreaWrite.Arg("file", "File to read data from").Required().StringVar(&pac.file)
	parameterAreaWrite.Flag("beginning-address", "address to start writing data to").Default("0x0000").Uint16Var(&pac.beginningAddress)

	parameterArea.Command("clear", "Clear parameter area").Action(pac.clearParameterArea)
}

func (pac *parameterAreaCommand) readParameterArea(c *kingpin.ParseContext) error {
	fmt.Printf("Reading from parameter area 0x%04x (%d words)\n", pac.beginningAddress, pac.readCount)
	var area fins.ParameterArea
	switch pac.area {
	case "plc":
		area = fins.ParameterAreaPLCSetup
	case "io":
		area = fins.ParameterAreaIOTableRegistration
	case "routing":
		area = fins.ParameterAreaRoutingTable
	case "bus":
		area = fins.ParameterAreaCPUBusUnitSetup
	case "unknown":
		area = fins.ParameterAreaUnknown
	}
	area, beginningAddress, readCount, data, e := (**pac.fins).ParameterAreaRead(area, pac.beginningAddress, pac.readCount)
	if e != nil {
		return e
	}

	var areaText string
	switch area {
	case fins.ParameterAreaPLCSetup:
		areaText = "PLC Setup"
	case fins.ParameterAreaIOTableRegistration:
		areaText = "I/O Table Registration"
	case fins.ParameterAreaRoutingTable:
		areaText = "Routing Table"
	case fins.ParameterAreaCPUBusUnitSetup:
		areaText = "CPU Bus Unit Setup"
	case fins.ParameterAreaUnknown:
		areaText = "Unknown"
	}

	fmt.Printf("Area: %s (0x%04x)\nBeginning Address: 0x%04x\nRead Count: %d\n", areaText, area, beginningAddress, readCount)
	fmt.Printf("%s", hex.Dump(data))

	return nil
}

func (pac *parameterAreaCommand) writeParameterArea(c *kingpin.ParseContext) error {
	fmt.Println("Writing from parameter area")
	return nil
}

func (pac *parameterAreaCommand) clearParameterArea(c *kingpin.ParseContext) error {
	fmt.Println("Clearing the parameter area")
	return nil
}
