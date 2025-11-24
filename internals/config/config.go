package config

import "time"

// AgentConfig holds all configuration values for the agent
type AgentConfig struct {
	ServerIP   string       `yaml:"server_ip"`
	ServerPort string       `yaml:"server_port"`
	Timing     TimingConfig `yaml:"timing"`
	Protocol   string       `yaml:"protocol"`
}

// ServerConfig holds all configuration values for the server
type ServerConfig struct {
	ListeningInterface string `yaml:"listening_interface"`
	ServerPort         string `yaml:"server_port"`
	Protocol           string `yaml:"protocol"`
	TlsKey             string `yaml:"tls_key"`
	TlsCert            string `yaml:"tls_cert"`
}

type TimingConfig struct {
	Delay  time.Duration `yaml:"delay"`
	Jitter int           `yaml:"jitter"`
}
