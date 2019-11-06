package fins

type header struct {
	messageType      byte
	responseRequired bool
	src              finsAddress
	dst              finsAddress
	serviceID        byte
	gatewayCount     byte
}

const (
	messageTypeCommand  byte = iota
	messageTypeResponse byte = iota
)

func defaultHeader(messageType byte, responseRequired bool, src finsAddress, dst finsAddress, serviceID byte) header {
	hdr := header{}
	hdr.messageType = messageType
	hdr.responseRequired = responseRequired
	hdr.gatewayCount = 2
	hdr.src = src
	hdr.dst = dst
	hdr.serviceID = serviceID
	return hdr
}

func defaultCommandHeader(src finsAddress, dst finsAddress, serviceID byte) header {
	hdr := defaultHeader(messageTypeCommand, true, src, dst, serviceID)
	return hdr
}

func defaultResponseHeader(commandHeader header) header {
	hdr := defaultHeader(messageTypeResponse, false, commandHeader.dst, commandHeader.src, commandHeader.serviceID)
	return hdr
}
