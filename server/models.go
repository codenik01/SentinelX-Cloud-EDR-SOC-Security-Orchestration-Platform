package main

import (
	"time"
)

// Host represents a registered endpoint agent.
type Host struct {
	HostID        string    `gorm:"type:varchar(100);primaryKey;column:host_id" json:"host_id"`
	Hostname      string    `gorm:"column:hostname;size:255" json:"hostname"`
	OS            string    `gorm:"column:os;size:100" json:"os"`
	Status        string    `gorm:"column:status;size:50" json:"status"` // "online", "offline"
	LastHeartbeat time.Time `gorm:"column:last_heartbeat" json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Event represents raw security telemetry logs ingested from the agents.
type Event struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HostID    string    `gorm:"type:varchar(100);index;column:host_id" json:"host_id"`
	Timestamp time.Time `gorm:"index" json:"timestamp"`
	Type      string    `gorm:"size:50" json:"type"` // "process", "network", "login", "file"
	Payload   string    `gorm:"type:text" json:"payload"` // JSON serialized event data
}

// Alert represents a triggered security event flag.
type Alert struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HostID      string    `gorm:"type:varchar(100);index;column:host_id" json:"host_id"`
	Host        Host      `gorm:"foreignKey:HostID;references:HostID" json:"host"`
	Title       string    `gorm:"size:255" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Severity    string    `gorm:"size:50" json:"severity"` // "critical", "high", "medium", "low"
	MatchedRule string    `gorm:"size:100" json:"matched_rule"`
	Payload     string    `gorm:"type:text" json:"payload"` // Raw event causing the alert
	Status      string    `gorm:"size:50;default:'unresolved'" json:"status"` // "unresolved", "resolved", "false_positive"
	Timestamp   time.Time `gorm:"index" json:"timestamp"`
}

// ThreatIntelItem represents known malicious indicators (IPs/Domains).
type ThreatIntelItem struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	IndicatorType string    `gorm:"size:50;index"` // "ip", "domain"
	Value         string    `gorm:"size:255;uniqueIndex"`
	Source        string    `gorm:"size:100"`
	Description   string    `gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
}
