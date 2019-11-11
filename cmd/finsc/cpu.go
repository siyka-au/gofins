package main

import (
	"fmt"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

func configureCPU(app *kingpin.Application, finsc *finscContext) {
	cpu := app.Command("cpu", "Remote CPU")
	cpu.Command("run", "Set the CPU operating mode to 'run'").Action(finsc.cpuOperatingModeRun)
	cpu.Command("monitor", "Set the CPU operating mode to 'monitor'").Action(finsc.cpuOperatingModeMonitor)
	cpu.Command("stop", "Set the CPU operating mode to 'stop'").Action(finsc.cpuOperatingModeStop)
	cpu.Command("read", "Read CPU unit information").Action(finsc.cpuUnitRead)
	cpu.Command("status", "Read CPU unit status").Action(finsc.cpuUnitStatus)
}

func (finsc *finscContext) cpuOperatingModeRun(c *kingpin.ParseContext) error {
	fmt.Println("Running CPU")
	return finsc.client.Run(fins.CPUOperatingModeRun)
}

func (finsc *finscContext) cpuOperatingModeMonitor(c *kingpin.ParseContext) error {
	fmt.Println("Monitoring CPU")
	return finsc.client.Run(fins.CPUOperatingModeMonitor)
}

func (finsc *finscContext) cpuOperatingModeStop(c *kingpin.ParseContext) error {
	fmt.Println("Stopping CPU")
	return finsc.client.Stop()
}

func (finsc *finscContext) cpuUnitRead(c *kingpin.ParseContext) error {
	fmt.Println("Reading CPU Unit")
	e := finsc.client.CPUUnitRead()
	return e
}

func (finsc *finscContext) cpuUnitStatus(c *kingpin.ParseContext) error {
	fmt.Println("CPU Unit Status")
	r, e := finsc.client.CPUUnitStatus()
	formatString := "%-31v %s\n"
	fmt.Printf(formatString, "Program status:", FormatBoolAs(r.Running, "running", "stopped"))
	fmt.Printf(formatString, "Flash memory writing:", FormatBoolAs(r.FlashMemoryWriting, "yes", "no"))
	fmt.Printf(formatString, "Battery present:", FormatBoolAs(r.BatteryPresent, "yes", "no"))
	fmt.Printf(formatString, "CPU status:", FormatBoolAs(r.Standby, "standby", "normal"))
	// Mode                      CPUOperatingMode
	fmt.Printf(formatString, "FALS error:", FormatBoolAs(r.FALSError, "yes", "no"))
	fmt.Printf(formatString, "Cycle time over:", FormatBoolAs(r.CycleTimeOver, "yes", "no"))
	fmt.Printf(formatString, "Program error:", FormatBoolAs(r.ProgramError, "yes", "no"))
	fmt.Printf(formatString, "I/O Setting error:", FormatBoolAs(r.IOSettingError, "yes", "no"))
	fmt.Printf(formatString, "I/O point overflow:", FormatBoolAs(r.IOPointOverflow, "yes", "no"))
	fmt.Printf(formatString, "Fatal inner board error:", FormatBoolAs(r.FatalInnerBoardError, "yes", "no"))
	fmt.Printf(formatString, "Duplication error:", FormatBoolAs(r.DuplicationError, "yes", "no"))
	fmt.Printf(formatString, "I/O bus error:", FormatBoolAs(r.IOBusError, "yes", "no"))
	fmt.Printf(formatString, "Memory error:", FormatBoolAs(r.MemoryError, "yes", "no"))
	fmt.Printf(formatString, "Other non-fatal error:", FormatBoolAs(r.OtherNonFatalError, "yes", "no"))
	fmt.Printf(formatString, "Special I/O unit setting error:", FormatBoolAs(r.SpecialIOUnitSettingError, "yes", "no"))
	fmt.Printf(formatString, "CS1 CPU bus unit setting error:", FormatBoolAs(r.CS1CPUBusUnitSettingError, "yes", "no"))
	fmt.Printf(formatString, "Battery error:", FormatBoolAs(r.BatteryError, "yes", "no"))
	fmt.Printf(formatString, "SYSMAC bus error:", FormatBoolAs(r.SYSMACBusError, "yes", "no"))
	fmt.Printf(formatString, "Special I/O unit error:", FormatBoolAs(r.SpecialIOUnitError, "yes", "no"))
	fmt.Printf(formatString, "CPU bus unit error:", FormatBoolAs(r.CPUBusUnitError, "yes", "no"))
	fmt.Printf(formatString, "Inner board error:", FormatBoolAs(r.InnerBoardError, "yes", "no"))
	fmt.Printf(formatString, "I/O verification error:", FormatBoolAs(r.IOVerificationError, "yes", "no"))
	fmt.Printf(formatString, "PLC setup error:", FormatBoolAs(r.PLCSetupError, "yes", "no"))
	fmt.Printf(formatString, "Basioc I/O unit error:", FormatBoolAs(r.BasicIOUnitError, "yes", "no"))
	fmt.Printf(formatString, "Interrupt task error:", FormatBoolAs(r.InterruptTaskError, "yes", "no"))
	fmt.Printf(formatString, "Duplex error:", FormatBoolAs(r.DuplexError, "yes", "no"))
	fmt.Printf(formatString, "FAL error:", FormatBoolAs(r.FALError, "yes", "no"))
	fmt.Printf(formatString, "FAL/FALS error code:", fmt.Sprintf("0x%04x", r.ErrorCode))
	fmt.Printf(formatString, "Message 0 present:", FormatBoolAs(r.MessagePresent[0], "yes", "no"))
	fmt.Printf(formatString, "Message 1 present:", FormatBoolAs(r.MessagePresent[1], "yes", "no"))
	fmt.Printf(formatString, "Message 2 present:", FormatBoolAs(r.MessagePresent[2], "yes", "no"))
	fmt.Printf(formatString, "Message 3 present:", FormatBoolAs(r.MessagePresent[3], "yes", "no"))
	fmt.Printf(formatString, "Message 4 present:", FormatBoolAs(r.MessagePresent[4], "yes", "no"))
	fmt.Printf(formatString, "Message 5 present:", FormatBoolAs(r.MessagePresent[5], "yes", "no"))
	fmt.Printf(formatString, "Message 6 present:", FormatBoolAs(r.MessagePresent[6], "yes", "no"))
	fmt.Printf(formatString, "Message 7 present:", FormatBoolAs(r.MessagePresent[7], "yes", "no"))
	return e
}
