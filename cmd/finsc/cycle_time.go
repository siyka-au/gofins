package main

import (
	"fmt"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cycleTimeCommand struct {
	fins   **fins.Client
}

func configureCycleTime(app *kingpin.Application, finsc *finscContext) {
	cycleTime := app.Command("cycle-time", "Cycle time statistics")
	ctc := &cycleTimeCommand{}
	ctc.fins = &finsc.client
	
	cycleTime.Command("initialise", "Initialise cycle time statistics").Action(ctc.cycleTimeInitialise)
	cycleTime.Command("read", "Read cycle time statistics").Action(ctc.cycleTimeRead)
}

func (ctc *cycleTimeCommand) cycleTimeInitialise(c *kingpin.ParseContext) error {
	fmt.Println("Initialising cycle time")
	e := (**ctc.fins).CycleTimeInitialise()
	return e
}

func (ctc *cycleTimeCommand) cycleTimeRead(c *kingpin.ParseContext) error {
	fmt.Println("Read cycle time statistics")
	avg, max, min, e :=(**ctc.fins).CycleTimeRead()
	if e != nil {
		return e
	}
	fmt.Printf("Avg: %s\n", avg)
	fmt.Printf("Max: %s\n", max)
	fmt.Printf("Min: %s\n", min)
	return nil
}
