package main

import (
	"flag"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
	"github.com/faanross/16_hr_course_rough_dev/internals/control"
	"github.com/faanross/16_hr_course_rough_dev/internals/server"
	"log"
	"os"
	"os/signal"
)

const pathToYAML = "./internals/config/config.yaml"

func main() {
	// Command line flag for config file path
	configPath := flag.String("config", pathToYAML, "path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadServerConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load our control API
	control.StartControlAPI()

	// Create BOTH servers regardless of config
	log.Printf("Starting both protocol servers on %s:%s", cfg.ListeningInterface, cfg.ServerPort)

	// Create HTTPS server
	httpsCfg := *cfg
	httpsCfg.Protocol = "https"
	httpsServer, err := server.NewServer(&httpsCfg)
	if err != nil {
		log.Fatalf("Failed to create HTTPS server: %v", err)
	}

	// Create DNS server
	dnsCfg := *cfg
	dnsCfg.Protocol = "dns"
	dnsServer, err := server.NewServer(&dnsCfg)
	if err != nil {
		log.Fatalf("Failed to create DNS server: %v", err)
	}

	// Start HTTPS server in goroutine
	go func() {
		log.Printf("Starting HTTPS server on %s:%s (TCP)", cfg.ListeningInterface, cfg.ServerPort)
		if err := httpsServer.Start(); err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	// Start DNS server in goroutine
	go func() {
		log.Printf("Starting DNS server on %s:%s (UDP)", cfg.ListeningInterface, cfg.ServerPort)
		if err := dnsServer.Start(); err != nil {
			log.Fatalf("DNS server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down both servers...")

	if err := httpsServer.Stop(); err != nil {
		log.Printf("Error HTTPS stopping server: %v", err)
	}

	if err := dnsServer.Stop(); err != nil {
		log.Printf("Error DNS stopping server: %v", err)
	}

}
