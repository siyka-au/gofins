package fins

import "fmt"

type EndCode uint16

// Data taken from Omron document Cat. No. W342-E1-15, pages 155-161
const (
	// EndCodeNormalCompletion End code: normal completion
	EndCodeNormalCompletion EndCode = 0x0000

	// EndCodeServiceInterrupted End code: normal completion; service was interrupted
	EndCodeServiceInterrupted EndCode = 0x0001

	// EndCodeLocalNodeNotInNetwork End code: local node error; local node not in network
	EndCodeLocalNodeNotInNetwork EndCode = 0x0101

	// EndCodeTokenTimeout End code: local node error; token timeout
	EndCodeTokenTimeout EndCode = 0x0102

	// EndCodeRetriesFailed End code: local node error; retries failed
	EndCodeRetriesFailed EndCode = 0x0103

	// EndCodeTooManySendFrames End code: local node error; too many send frames
	EndCodeTooManySendFrames EndCode = 0x0104

	// EndCodeNodeAddressRangeError End code: local node error; node address range error
	EndCodeNodeAddressRangeError EndCode = 0x0105

	// EndCodeNodeAddressRangeDuplication End code: local node error; node address range duplication
	EndCodeNodeAddressRangeDuplication EndCode = 0x0106

	// EndCodeDestinationNodeNotInNetwork End code: destination node error; destination node not in network
	EndCodeDestinationNodeNotInNetwork EndCode = 0x0201

	// EndCodeUnitMissing End code: destination node error; unit missing
	EndCodeUnitMissing EndCode = 0x0202

	// EndCodeThirdNodeMissing End code: destination node error; third node missing
	EndCodeThirdNodeMissing EndCode = 0x0203

	// EndCodeDestinationNodeBusy End code: destination node error; destination node busy
	EndCodeDestinationNodeBusy EndCode = 0x0204

	// EndCodeResponseTimeout End code: destination node error; response timeout
	EndCodeResponseTimeout EndCode = 0x0205

	// EndCodeCommunicationsControllerError End code: controller error; communication controller error
	EndCodeCommunicationsControllerError EndCode = 0x0301

	// EndCodeCPUUnitError End code: controller error; CPU unit error
	EndCodeCPUUnitError EndCode = 0x0302

	// EndCodeControllerError End code:  controller error; controller error
	EndCodeControllerError EndCode = 0x0303

	// EndCodeUnitNumberError End code: controller error; unit number error
	EndCodeUnitNumberError EndCode = 0x0304

	// EndCodeUndefinedCommand End code: service unsupported; undefined command
	EndCodeUndefinedCommand EndCode = 0x0401

	// EndCodeNotSupportedByModelVersion End code: service unsupported; not supported by model version
	EndCodeNotSupportedByModelVersion EndCode = 0x0402

	// EndCodeDestinationAddressSettingError End code: routing table error; destination address setting error
	EndCodeDestinationAddressSettingError EndCode = 0x0501

	// EndCodeNoRoutingTables End code: routing table error; no routing tables
	EndCodeNoRoutingTables EndCode = 0x0502

	// EndCodeRoutingTableError End code: routing table error; routing table error
	EndCodeRoutingTableError EndCode = 0x0503

	// EndCodeTooManyRelays End code: routing table error; too many relays
	EndCodeTooManyRelays EndCode = 0x0504

	// EndCodeCommandTooLong End code: command format error; command too long
	EndCodeCommandTooLong EndCode = 0x1001

	// EndCodeCommandTooShort End code: command format error; command too short
	EndCodeCommandTooShort EndCode = 0x1002

	// EndCodeElementsDataDontMatch End code: command format error; elements/data don't match
	EndCodeElementsDataDontMatch EndCode = 0x1003

	// EndCodeCommandFormatError End code: command format error; command format error
	EndCodeCommandFormatError EndCode = 0x1004

	// EndCodeHeaderError End code: command format error; header error
	EndCodeHeaderError EndCode = 0x1005

	// EndCodeAreaClassificationMissing End code: parameter error; classification missing
	EndCodeAreaClassificationMissing EndCode = 0x1101

	// EndCodeAccessSizeError End code: parameter error; access size error
	EndCodeAccessSizeError EndCode = 0x1102

	// EndCodeAddressRangeError End code: parameter error; address range error
	EndCodeAddressRangeError EndCode = 0x1103

	// EndCodeAddressRangeExceeded End code: parameter error; address range exceeded
	EndCodeAddressRangeExceeded EndCode = 0x1104

	// EndCodeProgramMissing End code: parameter error; program missing
	EndCodeProgramMissing EndCode = 0x1106

	// EndCodeRelationalError End code: parameter error; relational error
	EndCodeRelationalError EndCode = 0x1109

	// EndCodeDuplicateDataAccess End code: parameter error; duplicate data access
	EndCodeDuplicateDataAccess EndCode = 0x110a

	// EndCodeResponseTooBig End code: parameter error; response too big
	EndCodeResponseTooBig EndCode = 0x110b

	// EndCodeParameterError End code: parameter error
	EndCodeParameterError EndCode = 0x110c

	// EndCodeReadNotPossibleProtected End code: read not possible; protected
	EndCodeReadNotPossibleProtected EndCode = 0x2002

	// EndCodeReadNotPossibleTableMissing End code: read not possible; table missing
	EndCodeReadNotPossibleTableMissing EndCode = 0x2003

	// EndCodeReadNotPossibleDataMissing End code: read not possible; data missing
	EndCodeReadNotPossibleDataMissing EndCode = 0x2004

	// EndCodeReadNotPossibleProgramMissing End code: read not possible; program missing
	EndCodeReadNotPossibleProgramMissing EndCode = 0x2005

	// EndCodeReadNotPossibleFileMissing End code: read not possible; file missing
	EndCodeReadNotPossibleFileMissing EndCode = 0x2006

	// EndCodeReadNotPossibleDataMismatch End code: read not possible; data mismatch
	EndCodeReadNotPossibleDataMismatch EndCode = 0x2007

	// EndCodeWriteNotPossibleReadOnly End code: write not possible; read only
	EndCodeWriteNotPossibleReadOnly EndCode = 0x2101

	// EndCodeWriteNotPossibleProtected End code: write not possible; write protected
	EndCodeWriteNotPossibleProtected EndCode = 0x2102

	// EndCodeWriteNotPossibleCannotRegister End code: write not possible; cannot register
	EndCodeWriteNotPossibleCannotRegister EndCode = 0x2103

	// EndCodeWriteNotPossibleProgramMissing End code: write not possible; program missing
	EndCodeWriteNotPossibleProgramMissing EndCode = 0x2105

	// EndCodeWriteNotPossibleFileMissing End code: write not possible; file missing
	EndCodeWriteNotPossibleFileMissing EndCode = 0x2106

	// EndCodeWriteNotPossibleFileNameAlreadyExists End code: write not possible; file name already exists
	EndCodeWriteNotPossibleFileNameAlreadyExists EndCode = 0x2107

	// EndCodeWriteNotPossibleCannotChange End code: write not possible; cannot change
	EndCodeWriteNotPossibleCannotChange EndCode = 0x2108

	// EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution End code: not executeable in current mode during execution
	EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution EndCode = 0x2201

	// EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning End code: not executeable in current mode while running
	EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning EndCode = 0x2202

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram End code: not executeable in current mode; PLC is in PROGRAM mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram EndCode = 0x2203

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug End code: not executeable in current mode; PLC is in DEBUG mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug EndCode = 0x2204

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor End code: not executeable in current mode; PLC is in MONITOR mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor EndCode = 0x2205

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInRun End code: not executeable in current mode; PLC is in RUN mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInRun EndCode = 0x2206

	// EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode End code: not executeable in current mode; specified node is not polling node
	EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode EndCode = 0x2207

	// EndCodeNotExecutableInCurrentModeStepCannotBeExecuted End code: not executeable in current mode; step cannot be executed
	EndCodeNotExecutableInCurrentModeStepCannotBeExecuted EndCode = 0x2208

	// EndCodeNoSuchDeviceFileDeviceMissing End code: no such device; file device missing
	EndCodeNoSuchDeviceFileDeviceMissing EndCode = 0x2301

	// EndCodeNoSuchDeviceMemoryMissing End code: no such device; memory missing
	EndCodeNoSuchDeviceMemoryMissing EndCode = 0x2302

	// EndCodeNoSuchDeviceClockMissing End code: no such device; clock missing
	EndCodeNoSuchDeviceClockMissing EndCode = 0x2303

	// EndCodeCannotStartStopTableMissing End code: cannot start/stop; table missing
	EndCodeCannotStartStopTableMissing EndCode = 0x2401

	// EndCodeUnitErrorMemoryError End code: unit error; memory error
	EndCodeUnitErrorMemoryError EndCode = 0x2502

	// EndCodeUnitErrorIOError End code: unit error; IO error
	EndCodeUnitErrorIOError EndCode = 0x2503

	// EndCodeUnitErrorTooManyIOPoints End code: unit error; too many IO points
	EndCodeUnitErrorTooManyIOPoints EndCode = 0x2504

	// EndCodeUnitErrorCPUBusError End code: unit error; CPU bus error
	EndCodeUnitErrorCPUBusError EndCode = 0x2505

	// EndCodeUnitErrorIODuplication End code: unit error; IO duplication
	EndCodeUnitErrorIODuplication EndCode = 0x2506

	// EndCodeUnitErrorIOBusError End code: unit error; IO bus error
	EndCodeUnitErrorIOBusError EndCode = 0x2507

	// EndCodeUnitErrorSYSMACBUS2Error End code: unit error; SYSMAC BUS/2 error
	EndCodeUnitErrorSYSMACBUS2Error EndCode = 0x2509

	// EndCodeUnitErrorCPUBusUnitError End code: unit error; CPU bus unit error
	EndCodeUnitErrorCPUBusUnitError EndCode = 0x250a

	// EndCodeUnitErrorSYSMACBusNumberDuplication End code: unit error; SYSMAC bus number duplication
	EndCodeUnitErrorSYSMACBusNumberDuplication EndCode = 0x250d

	// EndCodeUnitErrorMemoryStatusError End code: unit error; memory status error
	EndCodeUnitErrorMemoryStatusError EndCode = 0x250f

	// EndCodeUnitErrorSYSMACBusTerminatorMissing End code: unit error; SYSMAC bus terminator missing
	EndCodeUnitErrorSYSMACBusTerminatorMissing EndCode = 0x2510

	// EndCodeCommandErrorNoProtection End code: command error; no protection
	EndCodeCommandErrorNoProtection EndCode = 0x2601

	// EndCodeCommandErrorIncorrectPassword End code: command error; incorrect password
	EndCodeCommandErrorIncorrectPassword EndCode = 0x2602

	// EndCodeCommandErrorProtected End code: command error; protected
	EndCodeCommandErrorProtected EndCode = 0x2604

	// EndCodeCommandErrorServiceAlreadyExecuting End code: command error; service already executing
	EndCodeCommandErrorServiceAlreadyExecuting EndCode = 0x2605

	// EndCodeCommandErrorServiceStopped End code: command error; service stopped
	EndCodeCommandErrorServiceStopped EndCode = 0x2606

	// EndCodeCommandErrorNoExecutionRight End code: command error; no execution right
	EndCodeCommandErrorNoExecutionRight EndCode = 0x2607

	// EndCodeCommandErrorSettingsNotComplete End code: command error; settings not complete
	EndCodeCommandErrorSettingsNotComplete EndCode = 0x2608

	// EndCodeCommandErrorNecessaryItemsNotSet End code: command error; necessary items not set
	EndCodeCommandErrorNecessaryItemsNotSet EndCode = 0x2609

	// EndCodeCommandErrorNumberAlreadyDefined End code: command error; number already defined
	EndCodeCommandErrorNumberAlreadyDefined EndCode = 0x260a

	// EndCodeCommandErrorErrorWillNotClear End code: command error; error will not clear
	EndCodeCommandErrorErrorWillNotClear EndCode = 0x260b

	// EndCodeAccessWriteErrorNoAccessRight End code: access write error; no access right
	EndCodeAccessWriteErrorNoAccessRight EndCode = 0x3001

	// EndCodeAbortServiceAborted End code: abort; service aborted
	EndCodeAbortServiceAborted EndCode = 0x4001
)

func (e EndCode) String() string {
	switch {
	case e == EndCodeNormalCompletion:
		return "normal completion"

	case e == EndCodeServiceInterrupted:
		return "normal completion; service was interrupted"

	case e == EndCodeLocalNodeNotInNetwork:
		return "local node error; local node not in network"

	case e == EndCodeTokenTimeout:
		return "local node error; token timeout"

	case e == EndCodeRetriesFailed:
		return "local node error; retries failed"

	case e == EndCodeTooManySendFrames:
		return "local node error; too many send frames"

	case e == EndCodeNodeAddressRangeError:
		return "local node error; node address range error"

	case e == EndCodeNodeAddressRangeDuplication:
		return "local node error; node address range duplication"

	case e == EndCodeDestinationNodeNotInNetwork:
		return "destination node error; destination node not in network"

	case e == EndCodeUnitMissing:
		return "destination node error; unit missing"

	case e == EndCodeThirdNodeMissing:
		return "destination node error; third node missing"

	case e == EndCodeDestinationNodeBusy:
		return "destination node error; destination node busy"

	case e == EndCodeResponseTimeout:
		return "destination node error; response timeout"

	case e == EndCodeCommunicationsControllerError:
		return "controller error; communication controller error"

	case e == EndCodeCPUUnitError:
		return "controller error; CPU unit error"

	case e == EndCodeControllerError:
		return " controller error; controller error"

	case e == EndCodeUnitNumberError:
		return "controller error; unit number error"

	case e == EndCodeUndefinedCommand:
		return "service unsupported; undefined command"

	case e == EndCodeNotSupportedByModelVersion:
		return "service unsupported; not supported by model version"

	case e == EndCodeDestinationAddressSettingError:
		return "routing table error; destination address setting error"

	case e == EndCodeNoRoutingTables:
		return "routing table error; no routing tables"

	case e == EndCodeRoutingTableError:
		return "routing table error; routing table error"

	case e == EndCodeTooManyRelays:
		return "routing table error; too many relays"

	case e == EndCodeCommandTooLong:
		return "command format error; command too long"

	case e == EndCodeCommandTooShort:
		return "command format error; command too short"

	case e == EndCodeElementsDataDontMatch:
		return "command format error; elements/data don't match"

	case e == EndCodeCommandFormatError:
		return "command format error; command format error"

	case e == EndCodeHeaderError:
		return "command format error; header error"

	case e == EndCodeAreaClassificationMissing:
		return "parameter error; classification missing"

	case e == EndCodeAccessSizeError:
		return "parameter error; access size error"

	case e == EndCodeAddressRangeError:
		return "parameter error; address range error"

	case e == EndCodeAddressRangeExceeded:
		return "parameter error; address range exceeded"

	case e == EndCodeProgramMissing:
		return "parameter error; program missing"

	case e == EndCodeRelationalError:
		return "parameter error; relational error"

	case e == EndCodeDuplicateDataAccess:
		return "parameter error; duplicate data access"

	case e == EndCodeResponseTooBig:
		return "parameter error; response too big"

	case e == EndCodeParameterError:
		return "parameter error"

	case e == EndCodeReadNotPossibleProtected:
		return "read not possible; protected"

	case e == EndCodeReadNotPossibleTableMissing:
		return "read not possible; table missing"

	case e == EndCodeReadNotPossibleDataMissing:
		return "read not possible; data missing"

	case e == EndCodeReadNotPossibleProgramMissing:
		return "read not possible; program missing"

	case e == EndCodeReadNotPossibleFileMissing:
		return "read not possible; file missing"

	case e == EndCodeReadNotPossibleDataMismatch:
		return "read not possible; data mismatch"

	case e == EndCodeWriteNotPossibleReadOnly:
		return "write not possible; read only"

	case e == EndCodeWriteNotPossibleProtected:
		return "write not possible; write protected"

	case e == EndCodeWriteNotPossibleCannotRegister:
		return "write not possible; cannot register"

	case e == EndCodeWriteNotPossibleProgramMissing:
		return "write not possible; program missing"

	case e == EndCodeWriteNotPossibleFileMissing:
		return "write not possible; file missing"

	case e == EndCodeWriteNotPossibleFileNameAlreadyExists:
		return "write not possible; file name already exists"

	case e == EndCodeWriteNotPossibleCannotChange:
		return "write not possible; cannot change"

	case e == EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution:
		return "not executeable in current mode during execution"

	case e == EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning:
		return "not executeable in current mode while running"

	case e == EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram:
		return "not executeable in current mode; PLC is in PROGRAM mode"

	case e == EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug:
		return "not executeable in current mode; PLC is in DEBUG mode"

	case e == EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor:
		return "not executeable in current mode; PLC is in MONITOR mode"

	case e == EndCodeNotExecutableInCurrentModeWrongPLCModeInRun:
		return "not executeable in current mode; PLC is in RUN mode"

	case e == EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode:
		return "not executeable in current mode; specified node is not polling node"

	case e == EndCodeNotExecutableInCurrentModeStepCannotBeExecuted:
		return "not executeable in current mode; step cannot be executed"

	case e == EndCodeNoSuchDeviceFileDeviceMissing:
		return "no such device; file device missing"

	case e == EndCodeNoSuchDeviceMemoryMissing:
		return "no such device; memory missing"

	case e == EndCodeNoSuchDeviceClockMissing:
		return "no such device; clock missing"

	case e == EndCodeCannotStartStopTableMissing:
		return "cannot start/stop; table missing"

	case e == EndCodeUnitErrorMemoryError:
		return "unit error; memory error"

	case e == EndCodeUnitErrorIOError:
		return "unit error; IO error"

	case e == EndCodeUnitErrorTooManyIOPoints:
		return "unit error; too many IO points"

	case e == EndCodeUnitErrorCPUBusError:
		return "unit error; CPU bus error"

	case e == EndCodeUnitErrorIODuplication:
		return "unit error; IO duplication"

	case e == EndCodeUnitErrorIOBusError:
		return "unit error; IO bus error"

	case e == EndCodeUnitErrorSYSMACBUS2Error:
		return "unit error; SYSMAC BUS/2 error"

	case e == EndCodeUnitErrorCPUBusUnitError:
		return "unit error; CPU bus unit error"

	case e == EndCodeUnitErrorSYSMACBusNumberDuplication:
		return "unit error; SYSMAC bus number duplication"

	case e == EndCodeUnitErrorMemoryStatusError:
		return "unit error; memory status error"

	case e == EndCodeUnitErrorSYSMACBusTerminatorMissing:
		return "unit error; SYSMAC bus terminator missing"

	case e == EndCodeCommandErrorNoProtection:
		return "command error; no protection"

	case e == EndCodeCommandErrorIncorrectPassword:
		return "command error; incorrect password"

	case e == EndCodeCommandErrorProtected:
		return "command error; protected"

	case e == EndCodeCommandErrorServiceAlreadyExecuting:
		return "command error; service already executing"

	case e == EndCodeCommandErrorServiceStopped:
		return "command error; service stopped"

	case e == EndCodeCommandErrorNoExecutionRight:
		return "command error; no execution right"

	case e == EndCodeCommandErrorSettingsNotComplete:
		return "command error; settings not complete"

	case e == EndCodeCommandErrorNecessaryItemsNotSet:
		return "command error; necessary items not set"

	case e == EndCodeCommandErrorNumberAlreadyDefined:
		return "command error; number already defined"

	case e == EndCodeCommandErrorErrorWillNotClear:
		return "command error; error will not clear"

	case e == EndCodeAccessWriteErrorNoAccessRight:
		return "access write error; no access right"

	case e == EndCodeAbortServiceAborted:
		return "abort; service aborted"

	}

	return fmt.Sprintf("unknown end code 0x%04x", uint16(e))
}
