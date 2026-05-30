package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password by default
		DB:       0,  // Standard DB index
	})

	log.Println("[SentinelX-Server] Connecting to Redis...")
	
	// Test ping connection
	err := RDB.Ping(Ctx).Err()
	if err != nil {
		log.Printf("[SentinelX-Server] Redis connection failed: %v. Running in localized memory fallback mode.", err)
		// We can still function or print, but in Docker Compose Redis will be available.
	} else {
		log.Println("[SentinelX-Server] Redis connection successful.")
	}
}

// TrackLoginAttempt increments a count in Redis for that host and source IP. Returns the current count.
func TrackLoginAttempt(hostID, sourceIP string, duration time.Duration) (int64, error) {
	if RDB == nil {
		return 0, fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf("brute:%s:%s", hostID, sourceIP)
	
	pipe := RDB.TxPipeline()
	incr := pipe.Incr(Ctx, key)
	pipe.Expire(Ctx, key, duration)
	
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

// ClearLoginAttempts flushes the count after an alert or successful mitigation
func ClearLoginAttempts(hostID, sourceIP string) {
	if RDB == nil {
		return
	}
	key := fmt.Sprintf("brute:%s:%s", hostID, sourceIP)
	RDB.Del(Ctx, key)
}
