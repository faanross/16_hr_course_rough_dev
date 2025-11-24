package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// LoadAgentConfig reads and parses the configuration file for the agent
func LoadAgentConfig(path string) (*AgentConfig, error) {

	// We'll provide path to *.yaml to function when we call it
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	// Create empty Config struct
	var agentCfg AgentConfig

	// Unmarshall YAML file values into struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&agentCfg); err != nil {
		return nil, fmt.Errorf("parsing agent config file: %w", err)
	}

	// Optional but good: Validate values
	if err := agentCfg.ValidateAgentConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &agentCfg, nil
}

// ValidateAgentConfig checks if the configuration is valid
func (ac *AgentConfig) ValidateAgentConfig() error {
	if ac.ServerIP == "" {
		return fmt.Errorf("server ip cannot be empty")
	}

	if ac.ServerPort == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if ac.Timing.Delay <= 0 {
		return fmt.Errorf("delay must be positive")
	}

	if ac.Timing.Jitter < 0 || ac.Timing.Jitter > 100 {
		return fmt.Errorf("jitter must be between 0 and 100")
	}

	if ac.Protocol != "dns" && ac.Protocol != "https" {
		return fmt.Errorf("protocol must be either 'dns' or 'https'")
	}

	return nil
}

// LoadServerConfig reads and parses the configuration file for the agent
func LoadServerConfig(path string) (*ServerConfig, error) {

	// We'll provide path to *.yaml to function when we call it
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	// Create empty Config struct
	var serverCfg ServerConfig

	// Unmarshall YAML file values into struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&serverCfg); err != nil {
		return nil, fmt.Errorf("parsing server config file: %w", err)
	}

	// Optional but good: Validate values
	if err := serverCfg.ValidateServerConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &serverCfg, nil
}

// ValidateServerConfig checks if the configuration is valid
func (sc *ServerConfig) ValidateServerConfig() error {

	if sc.ListeningInterface == "" {
		return fmt.Errorf("listening interface cannot be empty")
	}

	if sc.ServerPort == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if sc.Protocol != "dns" && sc.Protocol != "https" {
		return fmt.Errorf("protocol must be either 'dns' or 'https'")
	}

	if sc.TlsCert == "" || sc.TlsKey == "" {
		return fmt.Errorf("tls cert and key cannot be empty")
	}
	return nil
}
