package main

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIServerURL string
	HostID       string
	Hostname     string
	OS           string
	PollInterval time.Duration
	MockMode     bool
}

func LoadConfig() *Config {
	apiURL := os.Getenv("SENTINELX_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	hostID := os.Getenv("SENTINELX_HOST_ID")
	if hostID == "" {
		hostID = "host-01"
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "sentinelx-endpoint"
	}
	if envHost := os.Getenv("SENTINELX_HOSTNAME"); envHost != "" {
		hostname = envHost
	}

	osName := os.Getenv("SENTINELX_OS")
	if osName == "" {
		osName = "linux" // Default to linux target
	}

	intervalStr := os.Getenv("SENTINELX_POLL_INTERVAL_SEC")
	interval := 5 * time.Second
	if intervalStr != "" {
		if sec, err := strconv.Atoi(intervalStr); err == nil {
			interval = time.Duration(sec) * time.Second
		}
	}

	mockMode := true // Default to true to enable high-fidelity mock events out-of-the-box on non-Linux hosts
	if envMock := os.Getenv("SENTINELX_MOCK_MODE"); envMock == "false" {
		mockMode = false
	}

	return &Config{
		APIServerURL: apiURL,
		HostID:       hostID,
		Hostname:     hostname,
		OS:           osName,
		PollInterval: interval,
		MockMode:     mockMode,
	}
}
