package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Grpc        *grpc.GrpcConfig
	Kafka       *kafka.KafkaConfig
	RedisOption *redis.Options

	ExpirationInterval time.Duration
}

func InitConfig(ctx context.Context, log *jsonlog.Logger) (*Config, error) {
	cfg := &Config{}

	// Redis Configuration
	redisAddr := getEnv("REDIS_ADDR", "localhost:6378")
	log.PrintInfo(ctx, "Connecting to Redis...", map[string]string{"redis_address": redisAddr})
	cfg.RedisOption = &redis.Options{Addr: redisAddr}

	redisClient := redis.NewClient(cfg.RedisOption)
	defer redisClient.Close()

	// gRPC Configuration
	cfg.Grpc = &grpc.GrpcConfig{
        Port: getEnv("PORT", ":5055"),
		Host: getEnv("HOST", "localhost"),
	}

	// Expiration Interval
	intervalStr := getEnv("EXPIRATION_INTERVAL", "5m")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("invalid EXPIRATION_INTERVAL"), map[string]string{
			"interval": intervalStr,
		})
		return nil, fmt.Errorf("invalid EXPIRATION_INTERVAL: %v", err)
	}
	cfg.ExpirationInterval = interval

	// Kafka Configuration
	brokers := []string{getEnv("KAFKA_BROKER", "localhost:9092")}

	// Fetch Kafka topics
	topicsStr := getEnv("KAFKA_TOPICS", "[\"delete-expired-key-topic\", \"user-notifications\"]")
	if topicsStr == "" {
		log.PrintError(ctx, fmt.Errorf("KAFKA_TOPICS is empty"), nil)
		return nil, fmt.Errorf("KAFKA_TOPICS must be a non-empty JSON array")
	}

	var topics []string
	if err := json.Unmarshal([]byte(topicsStr), &topics); err != nil || len(topics) == 0 {
		log.PrintError(ctx, fmt.Errorf("failed to parse KAFKA_TOPICS"), map[string]string{
			"topics_raw": topicsStr,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("invalid or empty KAFKA_TOPICS: %v", err)
	}
	log.PrintInfo(ctx, fmt.Sprintf("Parsed Kafka Topics: %v", topics), nil)

	retriesStr := getEnv("KAFKA_RETRIES", "5")
	retries, err := strconv.Atoi(retriesStr)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("invalid KAFKA_RETRIES"), map[string]string{
			"retries_raw": retriesStr,
		})
		return nil, fmt.Errorf("invalid KAFKA_RETRIES: %v", err)
	}
	cfg.Kafka = &kafka.KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		MaxRetries: retries,
		GroupID:    "1",
	}


	return cfg, nil
}

// Helper function to fetch environment variables with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
