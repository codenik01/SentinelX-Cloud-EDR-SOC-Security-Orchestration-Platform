package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Local development fallback
		dsn = "host=localhost user=postgres password=postgres dbname=sentinelx port=5432 sslmode=disable"
	}

	var db *gorm.DB
	var err error

	// Retry database connection on startup to survive docker orchestration delay
	log.Println("[SentinelX-Server] Connecting to PostgreSQL...")
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("[SentinelX-Server] Database not ready, retrying in 3 seconds (%d/10)...", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("[SentinelX-Server] PostgreSQL connection failed: %v", err)
	}

	log.Println("[SentinelX-Server] Database connection successful. Running Auto-Migrations...")
	
	// Migrate models sequentially in dependency order
	err = db.AutoMigrate(&Host{})
	if err != nil {
		log.Fatalf("[SentinelX-Server] Host Migrate failed: %v", err)
	}
	err = db.AutoMigrate(&Event{})
	if err != nil {
		log.Fatalf("[SentinelX-Server] Event Migrate failed: %v", err)
	}
	err = db.AutoMigrate(&Alert{})
	if err != nil {
		log.Fatalf("[SentinelX-Server] Alert Migrate failed: %v", err)
	}
	err = db.AutoMigrate(&ThreatIntelItem{})
	if err != nil {
		log.Fatalf("[SentinelX-Server] ThreatIntel Migrate failed: %v", err)
	}

	DB = db
	log.Println("[SentinelX-Server] Database migrations complete.")

	SeedThreatIntel()
}

func SeedThreatIntel() {
	var count int64
	DB.Model(&ThreatIntelItem{}).Count(&count)
	if count > 0 {
		log.Println("[SentinelX-Server] Threat Intelligence database already seeded.")
		return
	}

	log.Println("[SentinelX-Server] Seeding Threat Intelligence Feed blocklist...")

	seeds := []ThreatIntelItem{
		{
			IndicatorType: "ip",
			Value:         "185.230.125.1",
			Source:        "Feodo Tracker",
			Description:   "Malicious C2 node linked to Ursnif/Gozi botnets",
			CreatedAt:     time.Now(),
		},
		{
			IndicatorType: "ip",
			Value:         "45.80.201.12",
			Source:        "BruteForceBlocker",
			Description:   "SSH brute forcing bot IP",
			CreatedAt:     time.Now(),
		},
		{
			IndicatorType: "domain",
			Value:         "evil-botnet.ru",
			Source:        "Spamhaus DBL",
			Description:   "Active phishing and credential harvesting control domain",
			CreatedAt:     time.Now(),
		},
		{
			IndicatorType: "domain",
			Value:         "hacker-c2-channel.net",
			Source:        "Emerging Threats",
			Description:   "Cobalt Strike command and control domain",
			CreatedAt:     time.Now(),
		},
	}

	for _, item := range seeds {
		if err := DB.Create(&item).Error; err != nil {
			log.Printf("[SentinelX-Server] Error seeding threat item %s: %v", item.Value, err)
		}
	}

	log.Println("[SentinelX-Server] Successfully seeded 4 threat intelligence indicators.")
}
