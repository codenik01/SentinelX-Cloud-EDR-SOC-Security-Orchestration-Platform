package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LinuxCollector struct {
	lastFileCheck time.Time
}

func NewLinuxCollector() *LinuxCollector {
	return &LinuxCollector{
		lastFileCheck: time.Now(),
	}
}

func (l *LinuxCollector) CollectProcesses() []ProcessEvent {
	var events []ProcessEvent
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return events
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		cmdBytes, err := ioutil.ReadFile(filepath.Join("/proc", file.Name(), "cmdline"))
		if err != nil {
			continue
		}
		command := string(cmdBytes)
		command = strings.ReplaceAll(command, "\x00", " ")
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		// Read UID
		statusBytes, err := ioutil.ReadFile(filepath.Join("/proc", file.Name(), "status"))
		if err != nil {
			continue
		}
		user := "unknown"
		scanner := bufio.NewScanner(strings.NewReader(string(statusBytes)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Uid:") {
				fields := strings.Fields(line)
				if len(fields) > 1 {
					uid := fields[1]
					if uid == "0" {
						user = "root"
					} else {
						user = "user-" + uid
					}
				}
				break
			}
		}

		events = append(events, ProcessEvent{
			PID:        pid,
			PPID:       1, // Simplified
			Command:    command,
			Executable: strings.Fields(command)[0],
			User:       user,
			Timestamp:  time.Now(),
		})
	}
	return events
}

func (l *LinuxCollector) CollectNetwork() []NetworkEvent {
	var events []NetworkEvent
	parseSockets := func(path string, proto string) {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return
		}
		lines := strings.Split(string(content), "\n")
		if len(lines) < 2 {
			return
		}
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			// Parse Local and Remote addresses (Format hex IP:Port)
			remoteHex := fields[2]
			parts := strings.Split(remoteHex, ":")
			if len(parts) != 2 {
				continue
			}
			ipHex, portHex := parts[0], parts[1]
			if ipHex == "00000000" { // Listening
				continue
			}
			// Convert hex IP to dot decimal
			ip, _ := parseHexIP(ipHex)
			port, _ := strconv.ParseInt(portHex, 16, 64)

			events = append(events, NetworkEvent{
				Source:      "127.0.0.1",
				Destination: ip,
				IP:          ip,
				Port:        int(port),
				Protocol:    proto,
				Process:     "kernel-socket",
				Timestamp:   time.Now(),
			})
		}
	}

	parseSockets("/proc/net/tcp", "TCP")
	parseSockets("/proc/net/udp", "UDP")
	return events
}

func (l *LinuxCollector) CollectLogins() []LoginEvent {
	var events []LoginEvent
	// Parse auth logs (e.g. /var/log/auth.log or /var/log/secure)
	paths := []string{"/var/log/auth.log", "/var/log/secure"}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		// Read last few lines to avoid huge ingestion overhead
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "sshd") && strings.Contains(line, "Accepted") {
				// E.g. "Accepted password for root from 192.168.1.100 port 55102 ssh2"
				fields := strings.Fields(line)
				user := "unknown"
				ip := "unknown"
				for i, field := range fields {
					if field == "for" && i+1 < len(fields) {
						user = fields[i+1]
					}
					if field == "from" && i+1 < len(fields) {
						ip = fields[i+1]
					}
				}
				events = append(events, LoginEvent{
					Event:     "login",
					User:      user,
					SourceIP:  ip,
					Timestamp: time.Now(),
				})
			} else if strings.Contains(line, "sshd") && strings.Contains(line, "Failed password") {
				fields := strings.Fields(line)
				user := "unknown"
				ip := "unknown"
				for i, field := range fields {
					if field == "for" && i+1 < len(fields) {
						user = fields[i+1]
					}
					if field == "from" && i+1 < len(fields) {
						ip = fields[i+1]
					}
				}
				events = append(events, LoginEvent{
					Event:     "failed_login",
					User:      user,
					SourceIP:  ip,
					Timestamp: time.Now(),
				})
			}
		}
	}
	return events
}

func (l *LinuxCollector) CollectFiles() []FileEvent {
	var events []FileEvent
	// Check /etc/passwd modification time
	info, err := os.Stat("/etc/passwd")
	if err == nil {
		modTime := info.ModTime()
		if modTime.After(l.lastFileCheck) {
			events = append(events, FileEvent{
				Path:      "/etc/passwd",
				Action:    "modify",
				User:      "root",
				Timestamp: modTime,
			})
			l.lastFileCheck = modTime
		}
	}
	return events
}

func parseHexIP(hexStr string) (string, error) {
	if len(hexStr) != 8 {
		return "", fmt.Errorf("invalid hex IP length")
	}
	d, _ := strconv.ParseInt(hexStr[6:8], 16, 64)
	c, _ := strconv.ParseInt(hexStr[4:6], 16, 64)
	b, _ := strconv.ParseInt(hexStr[2:4], 16, 64)
	a, _ := strconv.ParseInt(hexStr[0:2], 16, 64)
	return fmt.Sprintf("%d.%d.%d.%d", d, c, b, a), nil
}
