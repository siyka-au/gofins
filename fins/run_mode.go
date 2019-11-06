package fins

type CPUOperatingMode byte

const (
	MonitorMode CPUOperatingMode = 0x02
	RunMode     CPUOperatingMode = 0x04
)
