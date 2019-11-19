package fins

import (
	"encoding/binary"
)

// A FINS command request
type request struct {
	hdr         header
	commandCode uint16
	data        []byte
}

// A FINS command response
type response struct {
	hdr         header
	commandCode uint16
	endCode     uint16
	data        []byte
}

// A plc memory address to do a work
type memoryAddress struct {
	memoryArea MemoryArea
	address    uint16
	bitOffset  byte
}

func runCommand(opMode CPUOperatingMode) []byte {
	commandData := make([]byte, 2, 8)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeRun)
	commandData = append(commandData, []byte{0xff, 0xff}...)
	commandData = append(commandData, byte(opMode))
	return commandData
}

func stopCommand() []byte {
	commandData := make([]byte, 2, 4)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeStop)
	commandData = append(commandData, []byte{0xff, 0xff}...)
	return commandData
}

func cpuUnitStatusCommand() []byte {
	commandData := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeCPUUnitStatusRead)
	return commandData
}

func memoryAreaReadCommand(memoryAddr memoryAddress, itemCount uint16) []byte {
	commandData := make([]byte, 2, 8)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeMemoryAreaRead)
	commandData = append(commandData, encodeMemoryAddress(memoryAddr)...)
	commandData = append(commandData, []byte{0, 0}...)
	binary.BigEndian.PutUint16(commandData[6:8], itemCount)
	return commandData
}

func memoryAreaWriteCommand(memoryAddr memoryAddress, itemCount uint16, data []byte) []byte {
	commandData := make([]byte, 2, 8+len(data))
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeMemoryAreaWrite)
	commandData = append(commandData, encodeMemoryAddress(memoryAddr)...)
	commandData = append(commandData, []byte{0, 0}...)
	binary.BigEndian.PutUint16(commandData[6:8], itemCount)
	commandData = append(commandData, data...)
	return commandData
}

func memoryAreaFillCommand(memoryAddr memoryAddress, itemCount uint16, data byte) []byte {
	commandData := make([]byte, 2, 8+1)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeMemoryAreaFill)
	commandData = append(commandData, encodeMemoryAddress(memoryAddr)...)
	commandData = append(commandData, []byte{0, 0}...)
	binary.BigEndian.PutUint16(commandData[6:8], itemCount)
	commandData = append(commandData, data)
	return commandData
}

const (
	dayOfWeekSunday    byte = 0x00
	dayOfWeekMonday    byte = 0x01
	dayOfWeekTuesday   byte = 0x02
	dayOfWeekWednesday byte = 0x03
	dayOfWeekThursday  byte = 0x04
	dayOfWeekFriday    byte = 0x05
	dayOfWeekSaturday  byte = 0x06
)

func parameterAreaReadCommand(parameterArea ParameterArea, beginningWord uint16, numberOfWords uint16) []byte {
	commandData := make([]byte, 8, 8)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeParameterAreaRead)
	binary.BigEndian.PutUint16(commandData[2:4], uint16(parameterArea))
	binary.BigEndian.PutUint16(commandData[4:6], beginningWord)
	binary.BigEndian.PutUint16(commandData[6:8], numberOfWords & 0x3fff)

	return commandData
}

func clockReadCommand() []byte {
	commandData := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeClockRead)
	return commandData
}

func clockWriteCommand(year, month, day, hour, minute, second, dayOfWeek byte) []byte {
	commandData := make([]byte, 9, 9)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeClockWrite)
	commandData[2] = year
	commandData[3] = month
	commandData[4] = day
	commandData[5] = hour
	commandData[6] = minute
	commandData[7] = second
	commandData[8] = dayOfWeek
	return commandData
}

func cycleTimeCommand(subCommand byte) []byte {
	commandData := make([]byte, 3, 3)
	binary.BigEndian.PutUint16(commandData[0:2], CommandCodeCycleTimeRead)
	commandData[2] = subCommand
	return commandData
}

func cycleTimeInitialiseCommand() []byte {
	return cycleTimeCommand(0x00)
}

func cycleTimeReadCommand() []byte {
	return cycleTimeCommand(0x01)
}

func encodeMemoryAddress(memoryAddr memoryAddress) []byte {
	bytes := make([]byte, 4, 4)
	bytes[0] = byte(memoryAddr.memoryArea)
	binary.BigEndian.PutUint16(bytes[1:3], memoryAddr.address)
	bytes[3] = memoryAddr.bitOffset
	return bytes
}

func decodeMemoryAddress(data []byte) memoryAddress {
	return memoryAddress{MemoryArea(data[0]), binary.BigEndian.Uint16(data[1:3]), data[3]}
}

func decodeRequest(bytes []byte) request {
	return request{
		decodeHeader(bytes[0:10]),
		binary.BigEndian.Uint16(bytes[10:12]),
		bytes[12:],
	}
}

func encodeResponse(resp response) []byte {
	bytes := make([]byte, 4, 4+len(resp.data))
	binary.BigEndian.PutUint16(bytes[0:2], resp.commandCode)
	binary.BigEndian.PutUint16(bytes[2:4], resp.endCode)
	bytes = append(bytes, resp.data...)
	bh := encodeHeader(resp.hdr)
	bh = append(bh, bytes...)
	return bh
}

func decodeResponse(bytes []byte) response {
	return response{
		decodeHeader(bytes[0:10]),
		binary.BigEndian.Uint16(bytes[10:12]),
		binary.BigEndian.Uint16(bytes[12:14]),
		bytes[14:],
	}
}

const (
	icfBridgesBit          byte = 7
	icfMessageTypeBit      byte = 6
	icfResponseRequiredBit byte = 0
)

func encodeHeader(hdr header) []byte {
	var icf byte
	icf = 1 << icfBridgesBit
	if hdr.responseRequired == false {
		icf |= 1 << icfResponseRequiredBit
	}
	if hdr.messageType == messageTypeResponse {
		icf |= 1 << icfMessageTypeBit
	}
	bytes := []byte{
		icf,
		0x00,
		hdr.gatewayCount,
		hdr.dst.network, hdr.dst.node, hdr.dst.unit,
		hdr.src.network, hdr.src.node, hdr.src.unit,
		hdr.serviceID}
	return bytes
}

func decodeHeader(bytes []byte) header {
	hdr := header{}
	icf := bytes[0]
	if icf&1<<icfResponseRequiredBit == 0 {
		hdr.responseRequired = true
	}
	if icf&1<<icfMessageTypeBit == 0 {
		hdr.messageType = messageTypeCommand
	} else {
		hdr.messageType = messageTypeResponse
	}
	hdr.gatewayCount = bytes[2]
	hdr.dst = finsAddress{bytes[3], bytes[4], bytes[5]}
	hdr.src = finsAddress{bytes[6], bytes[7], bytes[8]}
	hdr.serviceID = bytes[9]

	return hdr
}
