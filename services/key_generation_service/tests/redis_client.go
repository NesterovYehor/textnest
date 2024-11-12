package testutils

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"golang.org/x/net/context"
)

// NewRedisClient creates a new Redis client for tests
func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Ensure Redis is running on this address
	})

	// Ensure Redis connection is available
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return nil
	}

	return client
}

