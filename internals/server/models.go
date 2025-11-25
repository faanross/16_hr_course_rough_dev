package server

import (
	"fmt"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
)

// NewServer creates a new server based on the protocol
func NewServer(cfg *config.ServerConfig) (Server, error) {
	switch cfg.Protocol {
	case "https":
		return NewHTTPSServer(cfg), nil
	case "dns":
		return nil, fmt.Errorf("DNS not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported protocol: %v", cfg.Protocol)
	}
}

// Server defines the contract for servers
type Server interface {
	// Start begins listening for requests
	Start() error

	// Stop gracefully shuts down the server
	Stop() error
}
