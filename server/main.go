package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[SentinelX-Server] Initializing SentinelX SOC Server Backend...")

	// 1. Initialize Postgres Database
	InitDB()

	// 2. Initialize Redis Cache (stubbed for now)
	InitRedis()

	// Start background housekeeper to mark stale heartbeats as offline
	go runAgentMonitorHousekeeper()

	// 3. Configure Gin Routing
	r := gin.Default()

	// Configure CORS for web client connection
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "sentinelx-api"})
	})

	// API endpoints
	api := r.Group("/api/v1")
	{
		api.POST("/heartbeat", HandleHeartbeat)
		api.POST("/events", HandleEvents)
		api.GET("/overview", GetOverview)
		api.GET("/alerts", GetAlerts)
		api.PUT("/alerts/:id/resolve", ResolveAlert)
		api.GET("/hosts", GetHosts)
		api.GET("/timeline", GetTimeline)
		api.POST("/threat-intel/sync", HandleThreatIntelSync)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[SentinelX-Server] SOC Backend API listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[SentinelX-Server] Gin runtime failed: %v", err)
	}
}

// Background worker checking for offline agents
func runAgentMonitorHousekeeper() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if DB == nil {
			continue
		}
		now := time.Now()
		// Mark hosts offline if no heartbeat in 15 seconds
		DB.Model(&Host{}).
			Where("status = ? AND last_heartbeat < ?", "online", now.Add(-15*time.Second)).
			Update("status", "offline")
	}
}

// Stubs removed - implemented in redis.go and detection.go
