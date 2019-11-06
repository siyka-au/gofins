package fins

// MemoryArea A FINS memory area
type MemoryArea byte

const (
	// MemoryAreaCIOBit Memory area: CIO area; bit
	MemoryAreaCIOBit MemoryArea = 0x30

	// MemoryAreaWRBit Memory area: work area; bit
	MemoryAreaWRBit MemoryArea = 0x31

	// MemoryAreaHRBit Memory area: holding area; bit
	MemoryAreaHRBit MemoryArea = 0x32

	// MemoryAreaARBit Memory area: axuillary area; bit
	MemoryAreaARBit MemoryArea = 0x33

	// MemoryAreaCIOWord Memory area: CIO area; word
	MemoryAreaCIOWord MemoryArea = 0xb0

	// MemoryAreaWRWord Memory area: work area; word
	MemoryAreaWRWord MemoryArea = 0xb1

	// MemoryAreaHRWord Memory area: holding area; word
	MemoryAreaHRWord MemoryArea = 0xb2

	// MemoryAreaARWord Memory area: auxillary area; word
	MemoryAreaARWord MemoryArea = 0xb3

	// MemoryAreaTimerCounterCompletionFlag Memory area: counter completion flag
	MemoryAreaTimerCounterCompletionFlag MemoryArea = 0x09

	// MemoryAreaTimerCounterPV Memory area: counter PV
	MemoryAreaTimerCounterPV MemoryArea = 0x89

	// MemoryAreaDMBit Memory area: data area; bit
	MemoryAreaDMBit MemoryArea = 0x02

	// MemoryAreaDMWord Memory area: data area; word
	MemoryAreaDMWord MemoryArea = 0x82

	// MemoryAreaTaskBit Memory area: task flags; bit
	MemoryAreaTaskBit MemoryArea = 0x06

	// MemoryAreaTaskStatus Memory area: task flags; status
	MemoryAreaTaskStatus MemoryArea = 0x46

	// MemoryAreaIndexRegisterPV Memory area: CIO bit
	MemoryAreaIndexRegisterPV MemoryArea = 0xdc

	// MemoryAreaDataRegisterPV Memory area: CIO bit
	MemoryAreaDataRegisterPV MemoryArea = 0xbc

	// MemoryAreaClockPulsesConditionFlagsBit Memory area: CIO bit
	MemoryAreaClockPulsesConditionFlagsBit MemoryArea = 0x07
)
