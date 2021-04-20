package fins

// Server Omron FINS server
type Server struct {
	provider ServerProvider
	addr     Address
}

// NewServer creates a new Omron FINS server
func NewServer(provider ServerProvider, addr Address) *Server {
	s := new(Server)
	s.provider = provider
	s.addr = addr

	return s
}

// CloseConnection Closes an Omron FINS connection
func (s *Server) Close() {
	s.provider.close()
}
