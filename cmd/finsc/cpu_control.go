package main

import (
	"fmt"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

func configureCpuControl(app *kingpin.Application, finsc *finscContext) {
	cpu := app.Command("cpu", "Remote CPU")
	cpu.Command("run", "Set the CPU operating mode to 'run'").Action(finsc.cpuControlRun)
	cpu.Command("monitor", "Set the CPU operating mode to 'monitor'").Action(finsc.cpuControlMonitor)
	cpu.Command("stop", "Set the CPU operating mode to 'stop'").Action(finsc.cpuControlStop)
}

func (finsc *finscContext) cpuControlRun(c *kingpin.ParseContext) error {
	fmt.Println("Running CPU")
	finsc.client.Run(fins.RunMode)
	return nil
}

func (finsc *finscContext) cpuControlMonitor(c *kingpin.ParseContext) error {
	fmt.Println("Monitoring CPU")
	finsc.client.Run(fins.MonitorMode)
	return nil
}

func (finsc *finscContext) cpuControlStop(c *kingpin.ParseContext) error {
	fmt.Println("Stopping CPU")
	finsc.client.Stop()
	return nil
}
