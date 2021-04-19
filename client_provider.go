package fins

// ClientProvider is the interface implements underlying methods.
type ClientProvider interface {

	// Close connection
	close() error

	// Send command
	sendCommand(header *Header, payload *Payload) (*Response, error)
}
