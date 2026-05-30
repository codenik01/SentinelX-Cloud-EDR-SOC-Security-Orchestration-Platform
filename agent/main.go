package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("[SentinelX-Agent] Starting SentinelX Endpoint Security Agent...")

	config := LoadConfig()
	log.Printf("[SentinelX-Agent] Target Server: %s", config.APIServerURL)
	log.Printf("[SentinelX-Agent] Endpoint Identity: %s (Hostname: %s)", config.HostID, config.Hostname)

	var collector Collector
	if config.MockMode {
		log.Println("[SentinelX-Agent] Running in DEMO/SIMULATION Telemetry Mode (Cross-Platform)...")
		collector = NewMockCollector()
	} else {
		log.Println("[SentinelX-Agent] Running in ACTIVE LINUX Telemetry Mode...")
		collector = NewLinuxCollector()
	}

	// Channel to signal graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Ingestion loop ticker
	ticker := time.NewTicker(config.PollInterval)
	defer ticker.Stop()

	// Initial Heartbeat / Registration
	sendHeartbeat(config)

	log.Printf("[SentinelX-Agent] Started collection daemon. Interval: %s", config.PollInterval)

	for {
		select {
		case <-ticker.C:
			// Collect telemetry
			payload := TelemetryPayload{
				HostID:        config.HostID,
				Hostname:      config.Hostname,
				OS:            config.OS,
				Timestamp:     time.Now(),
				ProcessEvents: collector.CollectProcesses(),
				NetworkEvents: collector.CollectNetwork(),
				LoginEvents:   collector.CollectLogins(),
				FileEvents:    collector.CollectFiles(),
			}

			// Only post if events were detected to avoid empty database clutter
			if len(payload.ProcessEvents) > 0 || len(payload.NetworkEvents) > 0 ||
				len(payload.LoginEvents) > 0 || len(payload.FileEvents) > 0 {
				
				log.Printf("[SentinelX-Agent] Capturing events: [Processes: %d, Network: %d, Logins: %d, Files: %d]",
					len(payload.ProcessEvents), len(payload.NetworkEvents), len(payload.LoginEvents), len(payload.FileEvents))
				
				sendTelemetry(config.APIServerURL+"/api/v1/events", payload)
			} else {
				// No telemetry, just send heartbeat
				sendHeartbeat(config)
			}

		case sig := <-stopChan:
			log.Printf("[SentinelX-Agent] Received shutdown signal (%v). Exiting cleanly...", sig)
			return
		}
	}
}

func sendHeartbeat(config *Config) {
	url := fmt.Sprintf("%s/api/v1/heartbeat", config.APIServerURL)
	
	heartbeatData := map[string]interface{}{
		"host_id":   config.HostID,
		"hostname":  config.Hostname,
		"os":        config.OS,
		"timestamp": time.Now(),
	}

	jsonData, err := json.Marshal(heartbeatData)
	if err != nil {
		log.Printf("[SentinelX-Agent] Error marshaling heartbeat: %v", err)
		return
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[SentinelX-Agent] Failed to send heartbeat to server: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[SentinelX-Agent] Heartbeat returned unexpected status: %d", resp.StatusCode)
	}
}

func sendTelemetry(url string, payload TelemetryPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[SentinelX-Agent] Error marshaling telemetry payload: %v", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[SentinelX-Agent] Error transmitting telemetry to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[SentinelX-Agent] Server returned error status code: %d", resp.StatusCode)
	}
}
