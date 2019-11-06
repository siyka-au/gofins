package main

import (
	"fmt"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

func configureCpuOperatingMode(app *kingpin.Application, finsc *finscContext) {
	cpuOperatingMode := app.Command("cpu", "Remote CPU operating mode")
	cpuOperatingMode.Command("run", "Set the CPU operating mode to 'run'").Action(finsc.cpuOperatingModeRun)
	cpuOperatingMode.Command("monitor", "Set the CPU operating mode to 'monitor'").Action(finsc.cpuOperatingModeMonitor)
	cpuOperatingMode.Command("stop", "Set the CPU operating mode to 'stop'").Action(finsc.cpuOperatingModeStop)
}

func (finsc *finscContext) cpuOperatingModeRun(c *kingpin.ParseContext) error {
	fmt.Println("Running CPU")
	finsc.client.Run(fins.RunMode)
	return nil
}

func (finsc *finscContext) cpuOperatingModeMonitor(c *kingpin.ParseContext) error {
	fmt.Println("Monitoring CPU")
	finsc.client.Run(fins.MonitorMode)
	return nil
}

func (finsc *finscContext) cpuOperatingModeStop(c *kingpin.ParseContext) error {
	fmt.Println("Stopping CPU")
	finsc.client.Stop()
	return nil
}
