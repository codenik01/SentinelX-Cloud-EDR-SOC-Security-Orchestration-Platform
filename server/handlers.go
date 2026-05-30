package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HeartbeatRequest binds the incoming agent status packet
type HeartbeatRequest struct {
	HostID    string    `json:"host_id" binding:"required"`
	Hostname  string    `json:"hostname" binding:"required"`
	OS        string    `json:"os" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type TelemetryPayload struct {
	HostID        string         `json:"host_id" binding:"required"`
	Hostname      string         `json:"hostname"`
	OS            string         `json:"os"`
	Timestamp     time.Time      `json:"timestamp"`
	ProcessEvents []ProcessEvent `json:"process_events"`
	NetworkEvents []NetworkEvent `json:"network_events"`
	LoginEvents   []LoginEvent   `json:"login_events"`
	FileEvents    []FileEvent    `json:"file_events"`
}

type ProcessEvent struct {
	PID        int       `json:"pid"`
	PPID       int       `json:"ppid"`
	Command    string    `json:"command"`
	Executable string    `json:"executable"`
	User       string    `json:"user"`
	Timestamp  time.Time `json:"timestamp"`
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
	Event     string    `json:"event"`
	User      string    `json:"user"`
	SourceIP  string    `json:"source_ip"`
	Timestamp time.Time `json:"timestamp"`
}

type FileEvent struct {
	Path      string    `json:"path"`
	Action    string    `json:"action"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

// HandleHeartbeat updates or inserts endpoint records
func HandleHeartbeat(c *gin.Context) {
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var host Host
	result := DB.Where("host_id = ?", req.HostID).First(&host)

	if result.Error != nil {
		// Create new host record
		host = Host{
			HostID:        req.HostID,
			Hostname:      req.Hostname,
			OS:            req.OS,
			Status:        "online",
			LastHeartbeat: time.Now(),
		}
		DB.Create(&host)
	} else {
		// Update existing host status
		host.Hostname = req.Hostname
		host.OS = req.OS
		host.Status = "online"
		host.LastHeartbeat = time.Now()
		DB.Save(&host)
	}

	c.JSON(http.StatusOK, gin.H{"status": "registered", "host_id": host.HostID})
}

// HandleEvents records agent telemetry and passes to detection engine
func HandleEvents(c *gin.Context) {
	var payload TelemetryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update host last seen
	DB.Model(&Host{}).Where("host_id = ?", payload.HostID).Updates(map[string]interface{}{
		"status":         "online",
		"last_heartbeat": time.Now(),
	})

	// Process event batches
	go processTelemetryAsync(payload)

	c.JSON(http.StatusCreated, gin.H{"status": "accepted", "records": len(payload.ProcessEvents) + len(payload.NetworkEvents) + len(payload.LoginEvents) + len(payload.FileEvents)})
}

func processTelemetryAsync(payload TelemetryPayload) {
	// Save Processes
	for _, pe := range payload.ProcessEvents {
		data, _ := json.Marshal(pe)
		evt := Event{
			HostID:    payload.HostID,
			Timestamp: pe.Timestamp,
			Type:      "process",
			Payload:   string(data),
		}
		DB.Create(&evt)
		
		// Detection Engine Hook (added in next commit)
		AnalyzeProcessEvent(payload.HostID, pe, evt.ID)
	}

	// Save Network Sockets
	for _, ne := range payload.NetworkEvents {
		data, _ := json.Marshal(ne)
		evt := Event{
			HostID:    payload.HostID,
			Timestamp: ne.Timestamp,
			Type:      "network",
			Payload:   string(data),
		}
		DB.Create(&evt)

		AnalyzeNetworkEvent(payload.HostID, ne, evt.ID)
	}

	// Save Login attempts
	for _, le := range payload.LoginEvents {
		data, _ := json.Marshal(le)
		evt := Event{
			HostID:    payload.HostID,
			Timestamp: le.Timestamp,
			Type:      "login",
			Payload:   string(data),
		}
		DB.Create(&evt)

		AnalyzeLoginEvent(payload.HostID, le, evt.ID)
	}

	// Save File Integrity events
	for _, fe := range payload.FileEvents {
		data, _ := json.Marshal(fe)
		evt := Event{
			HostID:    payload.HostID,
			Timestamp: fe.Timestamp,
			Type:      "file",
			Payload:   string(data),
		}
		DB.Create(&evt)

		AnalyzeFileEvent(payload.HostID, fe, evt.ID)
	}
}

// GetOverview calculates the live dashboard gauges
func GetOverview(c *gin.Context) {
	var totalHosts int64
	var onlineHosts int64
	var totalEvents int64
	var totalAlerts int64

	DB.Model(&Host{}).Count(&totalHosts)
	DB.Model(&Host{}).Where("status = ?", "online").Count(&onlineHosts)
	DB.Model(&Event{}).Count(&totalEvents)
	DB.Model(&Alert{}).Count(&totalAlerts)

	var severities []struct {
		Severity string
		Count    int64
	}
	DB.Model(&Alert{}).Select("severity, count(*) as count").Group("severity").Scan(&severities)

	sevMap := map[string]int64{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}
	for _, s := range severities {
		sevMap[s.Severity] = s.Count
	}

	// Last 24 hours alert trends
	var trend []struct {
		Time  time.Time
		Count int64
	}
	DB.Model(&Alert{}).Select("date_trunc('hour', timestamp) as time, count(*) as count").
		Where("timestamp > ?", time.Now().Add(-24*time.Hour)).
		Group("date_trunc('hour', timestamp)").
		Order("time ASC").
		Scan(&trend)

	c.JSON(http.StatusOK, gin.H{
		"total_hosts":     totalHosts,
		"online_hosts":    onlineHosts,
		"total_events":    totalEvents,
		"total_alerts":    totalAlerts,
		"critical_alerts": sevMap["critical"],
		"high_alerts":     sevMap["high"],
		"medium_alerts":   sevMap["medium"],
		"low_alerts":      sevMap["low"],
		"alert_trend":     trend,
	})
}

// GetAlerts queries security notifications
func GetAlerts(c *gin.Context) {
	var alerts []Alert
	query := DB.Preload("Host").Order("timestamp DESC")

	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if hostID := c.Query("host_id"); hostID != "" {
		query = query.Where("host_id = ?", hostID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	query = query.Limit(limit)

	if err := query.Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// ResolveAlert updates alert ticket status
func ResolveAlert(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"` // "resolved", "false_positive"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Model(&Alert{}).Where("id = ?", id).Update("status", req.Status)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// GetHosts lists monitored host configurations
func GetHosts(c *gin.Context) {
	var hosts []Host
	if err := DB.Order("hostname ASC").Find(&hosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Dynamic online check (offline if heartbeat > 15s ago)
	now := time.Now()
	for i := range hosts {
		if now.Sub(hosts[i].LastHeartbeat) > 15*time.Second && hosts[i].Status == "online" {
			hosts[i].Status = "offline"
			DB.Model(&Host{}).Where("host_id = ?", hosts[i].HostID).Update("status", "offline")
		}
	}

	c.JSON(http.StatusOK, hosts)
}

// GetTimeline compiles events in reverse chronological order
func GetTimeline(c *gin.Context) {
	var events []Event
	limit := 50
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil {
			limit = l
		}
	}

	if err := DB.Order("timestamp DESC").Limit(limit).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// HandleThreatIntelSync adds external threat intel feed data
func HandleThreatIntelSync(c *gin.Context) {
	var req []struct {
		IndicatorType string `json:"type" binding:"required"`
		Value         string `json:"value" binding:"required"`
		Source        string `json:"source"`
		Description   string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedCount := 0
	for _, item := range req {
		intel := ThreatIntelItem{
			IndicatorType: item.IndicatorType,
			Value:         item.Value,
			Source:        item.Source,
			Description:   item.Description,
			CreatedAt:     time.Now(),
		}
		// Try to create, skip if unique constraint fails
		if err := DB.Create(&intel).Error; err == nil {
			addedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "synchronized",
		"synced_items": len(req),
		"new_items":    addedCount,
	})
}
