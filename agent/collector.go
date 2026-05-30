package main

import "time"

type ProcessEvent struct {
	PID         int       `json:"pid"`
	PPID        int       `json:"ppid"`
	Command     string    `json:"command"`
	Executable  string    `json:"executable"`
	User        string    `json:"user"`
	Timestamp   time.Time `json:"timestamp"`
}

type NetworkEvent struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`
	Process     string    `json:"process"`
	Timestamp   time.Time `json:"timestamp"`
}

type LoginEvent struct {
	Event     string    `json:"event"` // "login" or "failed_login"
	User      string    `json:"user"`
	SourceIP  string    `json:"source_ip"`
	Timestamp time.Time `json:"timestamp"`
}

type FileEvent struct {
	Path      string    `json:"path"`
	Action    string    `json:"action"` // "create", "modify", "delete"
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

type TelemetryPayload struct {
	HostID         string          `json:"host_id"`
	Hostname       string          `json:"hostname"`
	OS             string          `json:"os"`
	Timestamp      time.Time       `json:"timestamp"`
	ProcessEvents  []ProcessEvent  `json:"process_events,omitempty"`
	NetworkEvents  []NetworkEvent  `json:"network_events,omitempty"`
	LoginEvents    []LoginEvent    `json:"login_events,omitempty"`
	FileEvents     []FileEvent     `json:"file_events,omitempty"`
}

type Collector interface {
	CollectProcesses() []ProcessEvent
	CollectNetwork() []NetworkEvent
	CollectLogins() []LoginEvent
	CollectFiles() []FileEvent
}
