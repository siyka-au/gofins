package fins

// ServerProvider is the interface implements underlying methods.
type ServerProvider interface {
	// Send command
	// register(handler func(header *Header, payload *Payload) (*Header, *Payload)) error

	// Close connection
	close() error
}
