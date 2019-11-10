package fins

type CPUOperatingMode byte

const (
	CPUOperatingModeStopProgram CPUOperatingMode = 0x00
	CPUOperatingModeMonitor     CPUOperatingMode = 0x02
	CPUOperatingModeRun         CPUOperatingMode = 0x04
)
