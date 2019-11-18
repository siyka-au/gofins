package main

import (
	"fmt"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

type parameterAreaCommand struct {
	fins   **fins.Client
	beginningAddress uint32
	file	string
}

func configureParameterArea(app *kingpin.Application, finsc *finscContext) {
	parameterArea := app.Command("parameter-area", "Parameter area")
	pac := &parameterAreaCommand{}
	pac.fins = &finsc.client
	
	parameterAreaRead := parameterArea.Command("read", "Read parameter area").Action(pac.readParameterArea)
	parameterAreaRead.Arg("file", "File to save data to").Required().StringVar(&pac.file)
	parameterAreaRead.Flag("beginning-address", "address to start reading data from").Default("0x00000000").Uint32Var(&pac.beginningAddress)

	parameterAreaWrite := parameterArea.Command("write", "Write parameter area").Action(pac.writeParameterArea)
	parameterAreaWrite.Arg("file", "File to read data from").Required().StringVar(&pac.file)
	parameterAreaWrite.Flag("beginning-address", "address to start writing data to").Default("0x00000000").Uint32Var(&pac.beginningAddress)

	parameterArea.Command("clear", "Clear parameter area").Action(pac.clearParameterArea)
}

func (pac *parameterAreaCommand) readParameterArea(c *kingpin.ParseContext) error {
	fmt.Println("Reading from parameter area")
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
