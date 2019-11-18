package fins

// ParameterArea A FINS parameter area
type ParameterArea uint16

const (
	// ParameterAreaPLCSetup PLC setup area; total 512 words
	ParameterAreaPLCSetup ParameterArea = 0x8010

	// ParameterAreaIOTableRegistration I/O registration table area; total 1280 words
	ParameterAreaIOTableRegistration ParameterArea = 0x8012

	// ParameterAreaRoutingTable Routing table area; total 512 words
	ParameterAreaRoutingTable ParameterArea = 0x8013

	// ParameterAreaCPUBusUnitSetup CPU bus unit setup area; total 5184 words
	ParameterAreaCPUBusUnitSetup ParameterArea = 0x8002

	// ParameterAreaUnknown Memory area: CIO area; word
	ParameterAreaUnknown ParameterArea = 0x8000
)
