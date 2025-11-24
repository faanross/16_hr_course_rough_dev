package main

import (
	"flag"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
	"log"
)

const pathToYAML = "./configs/config.yaml"

func main() {
	// Command line flag for config file path
	configPath := flag.String("config", pathToYAML, "path to configuration file")
	flag.Parse()

	// Load configuration
	agentCfg, err := config.LoadAgentConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Loaded configuration:\n")
	log.Printf("-> Server IP: %s\n", agentCfg.ServerIP)
	log.Printf("-> Server Port: %s\n", agentCfg.ServerPort)
	log.Printf("-> Delay: %v\n", agentCfg.Timing.Delay)
	log.Printf("-> Jitter: %d%%\n", agentCfg.Timing.Jitter)
	log.Printf("-> Starting Protocol: %s\n", agentCfg.Protocol)

}
