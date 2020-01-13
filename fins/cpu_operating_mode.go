package fins

// CPUOperatingMode The operating mode of the CPU
type CPUOperatingMode byte

const (
	// CPUOperatingModeStopProgram CPU is stopped
	CPUOperatingModeStopProgram CPUOperatingMode = 0x00

	// CPUOperatingModeMonitor CPU is monitoring
	CPUOperatingModeMonitor     CPUOperatingMode = 0x02

	// CPUOperatingModeRun CPU is running
	CPUOperatingModeRun         CPUOperatingMode = 0x04
)
