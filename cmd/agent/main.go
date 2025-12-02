package main

import (
	"context"
	"github.com/faanross/16_hr_course_rough_dev/internals/agent"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
	"log"
	"os"
	"os/signal"
)

func main() {

	// Load configuration from embedded bytes
	cfg, err := config.LoadAgentConfigFromBytes(config.EmbeddedAgentConfig)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Call our factory function
	comm, err := agent.NewAgent(cfg)
	if err != nil {
		log.Fatalf("Failed to create communicator: %v", err)
	}

	// ALL THIS DOWN HERE IS THE NEW CODE
	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Printf("Starting %s client run loop", cfg.Protocol)
		log.Printf("Delay: %v, Jitter: %d%%", cfg.Timing.Delay, cfg.Timing.Jitter)

		if err := agent.RunLoop(ctx, comm, cfg); err != nil {
			log.Printf("Run loop error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Shutting down client...")
	cancel() // This will cause the run loop to exit
}
