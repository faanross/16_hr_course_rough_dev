package config

import (
	_ "embed"
)

//go:embed config.yaml
var EmbeddedAgentConfig []byte
