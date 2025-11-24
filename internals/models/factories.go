package models

import (
	"fmt"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
)

// NewAgent creates a new communicator based on the protocol
func NewAgent(cfg *config.AgentConfig) (Agent, error) {
	switch cfg.Protocol {
	case "https":
		return nil, fmt.Errorf("HTTPS not yet implemented")
	case "dns":
		return nil, fmt.Errorf("DNS not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported protocol: %v", cfg.Protocol)
	}
}

// NewServer creates a new server based on the protocol
func NewServer(cfg *config.ServerConfig) (Server, error) {
	switch cfg.Protocol {
	case "https":
		return nil, fmt.Errorf("HTTPS not yet implemented")
	case "dns":
		return nil, fmt.Errorf("DNS not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported protocol: %v", cfg.Protocol)
	}
}
