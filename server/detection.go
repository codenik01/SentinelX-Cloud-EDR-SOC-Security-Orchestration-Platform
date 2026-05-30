package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// AnalyzeProcessEvent evaluates process executions against EDR heuristics
func AnalyzeProcessEvent(hostID string, pe ProcessEvent, evtID uint) {
	cmd := strings.ToLower(pe.Command)

	// Rule 1: Reverse Shell Detection
	isReverseShell := false
	reason := ""
	if strings.Contains(cmd, "nc ") && (strings.Contains(cmd, "bash") || strings.Contains(cmd, "sh ") || strings.Contains(cmd, "-e ")) {
		isReverseShell = true
		reason = "Netcat spawning shell command observed: " + pe.Command
	} else if strings.Contains(cmd, "/dev/tcp/") && (strings.Contains(cmd, "bash -i") || strings.Contains(cmd, "sh -i")) {
		isReverseShell = true
		reason = "Bash redirection socket observed: " + pe.Command
	} else if strings.Contains(cmd, "python") && strings.Contains(cmd, "pty.spawn") {
		isReverseShell = true
		reason = "Python interactive pty allocation observed: " + pe.Command
	}

	if isReverseShell {
		createAlert(hostID, "Reverse Shell Detected", reason, "critical", "SIG-REVERSE-SHELL", pe, pe.Timestamp)
	}

	// Rule 2: Suspicious User/Admin Creation Heuristic
	isUserMod := false
	if strings.Contains(cmd, "useradd ") || strings.Contains(cmd, "usermod ") || strings.Contains(cmd, "groupadd ") || strings.Contains(cmd, "chsh ") {
		isUserMod = true
		reason = "System administrative command executed: " + pe.Command
	}

	if isUserMod {
		createAlert(hostID, "Suspicious User/Admin Modification", reason, "high", "SIG-USER-MOD", pe, pe.Timestamp)
	}
}

// AnalyzeNetworkEvent inspects outbound connection mappings
func AnalyzeNetworkEvent(hostID string, ne NetworkEvent, evtID uint) {
	// Look up Destination IP/Domain against Threat Intel Database
	var count int64
	var item ThreatIntelItem
	
	// Check IP
	DB.Model(&ThreatIntelItem{}).Where("indicator_type = ? AND value = ?", "ip", ne.IP).First(&item).Count(&count)
	if count > 0 {
		reason := fmt.Sprintf("Outbound network connection to blacklisted IP: %s (Source: %s, Desc: %s)", ne.IP, item.Source, item.Description)
		createAlert(hostID, "Threat Intelligence Match (IP)", reason, "critical", "SIG-INTEL-IP", ne, ne.Timestamp)
		return
	}

	// Check Domain
	if ne.Destination != "" && ne.Destination != ne.IP {
		DB.Model(&ThreatIntelItem{}).Where("indicator_type = ? AND value = ?", "domain", strings.ToLower(ne.Destination)).First(&item).Count(&count)
		if count > 0 {
			reason := fmt.Sprintf("Outbound network connection to blacklisted domain: %s (Source: %s, Desc: %s)", ne.Destination, item.Source, item.Description)
			createAlert(hostID, "Threat Intelligence Match (Domain)", reason, "critical", "SIG-INTEL-DOMAIN", ne, ne.Timestamp)
			return
		}
	}
}

// AnalyzeLoginEvent processes authentication attempts
func AnalyzeLoginEvent(hostID string, le LoginEvent, evtID uint) {
	if le.Event == "failed_login" {
		// Stateful Rule: SSH Brute Force (5 failures in 1 minute using Redis counters)
		count, err := TrackLoginAttempt(hostID, le.SourceIP, 1*time.Minute)
		if err != nil {
			log.Printf("[Detection-Engine] Redis tracking error: %v", err)
			return
		}

		log.Printf("[Detection-Engine] Host: %s, Attacker: %s, Failed attempts: %d/5", hostID, le.SourceIP, count)

		if count == 5 { // Fire alert exactly at 5 to prevent spamming
			reason := fmt.Sprintf("Host has experienced 5 failed logins within 1 minute from source IP: %s", le.SourceIP)
			createAlert(hostID, "SSH Brute Force Attempt", reason, "medium", "SIG-SSH-BRUTE", le, le.Timestamp)
		}
	} else if le.Event == "login" {
		// Log logins, clear counters
		ClearLoginAttempts(hostID, le.SourceIP)
		
		// If user is root, generate high-security notification
		if le.User == "root" {
			reason := fmt.Sprintf("Successful login to root shell from source IP: %s", le.SourceIP)
			createAlert(hostID, "Root Login Detected", reason, "high", "SIG-ROOT-LOGIN", le, le.Timestamp)
		}
	}
}

// AnalyzeFileEvent checks for changes in sensitive files
func AnalyzeFileEvent(hostID string, fe FileEvent, evtID uint) {
	path := strings.ToLower(fe.Path)

	// Rule: Sensitive credentials file modification
	if path == "/etc/passwd" || path == "/etc/shadow" || path == "/etc/gshadow" {
		reason := fmt.Sprintf("Sensitive authentication registry file write: %s (User: %s)", fe.Path, fe.User)
		createAlert(hostID, "Sensitive System File Modified", reason, "high", "SIG-FILE-PASSWD", fe, fe.Timestamp)
	}

	// Rule: Persistence modification via cron
	if strings.Contains(path, "cron") {
		reason := fmt.Sprintf("Cron configuration modified for persistence: %s (User: %s)", fe.Path, fe.User)
		createAlert(hostID, "New Cron Job Added", reason, "low", "SIG-FILE-CRON", fe, fe.Timestamp)
	}
}

// Helper to create and record alert into Postgres
func createAlert(hostID, title, description, severity, rule string, payload interface{}, timestamp time.Time) {
	payloadBytes, _ := json.Marshal(payload)

	alert := Alert{
		HostID:      hostID,
		Title:       title,
		Description: description,
		Severity:    severity,
		MatchedRule: rule,
		Payload:     string(payloadBytes),
		Status:      "unresolved",
		Timestamp:   timestamp,
	}

	if err := DB.Create(&alert).Error; err != nil {
		log.Printf("[Detection-Engine] Failed to persist security alert: %v", err)
	} else {
		log.Printf("[Detection-Engine] !!! ALERT CREATED [%s] !!! Host: %s, Rule: %s, Desc: %s", 
			strings.ToUpper(severity), hostID, rule, description)
	}
}
