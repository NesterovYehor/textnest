package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
)

type Config struct {
	Grpc        *grpc.GrpcConfig
	Kafka       *kafka.KafkaConfig
	RedisOption *redis.Options

	ExpirationInterval time.Duration
}

func InitConfig(ctx context.Context, log *jsonlog.Logger) (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.PrintError(ctx, fmt.Errorf("Warning: .env file not found, relying on environment variables"), nil)
		return nil, err
	}

	cfg := &Config{}

	// Redis Configuration
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	cfg.RedisOption = &redis.Options{Addr: redisAddr}
	if _, err := redis.NewClient(cfg.RedisOption).Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	// gRPC Configuration
	cfg.Grpc = &grpc.GrpcConfig{
		Port: getEnv("PORT", "5555"),
		Host: getEnv("HOST", "localhost"),
	}

	// Expiration Interval
	intervalStr := getEnv("EXPIRATION_INTERVAL", "5m")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid EXPIRATION_INTERVAL: %w", err)
	}
	cfg.ExpirationInterval = interval

	// Kafka Configuration
	brokers := []string{getEnv("KAFKA_BROKER", "localhost:9092")}

	// Fetch Kafka topics
	kafkaTopics := os.Getenv("KAFKA_TOPICS")
	var topics []string
	if kafkaTopics != "" {
		err := json.Unmarshal([]byte(kafkaTopics), &topics)
		if err != nil {
			log.PrintError(ctx, fmt.Errorf("failed to unmarshal KAFKA_TOPICS: %w", err), nil)
			return nil, fmt.Errorf("failed to unmarshal KAFKA_TOPICS: %w", err)
		}
		fmt.Println("KAFKA_TOPICS unmarshalled successfully:", topics)

	} else {
		// Default to empty map if no topics are specified
		topics = make([]string, 0, 1)
	}

	retries, err := strconv.Atoi(getEnv("KAFKA_RETRIES", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid KAFKA_RETRIES: %w", err)
	}
	cfg.Kafka = &kafka.KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		MaxRetries: retries,
		GroupID:    "1",
	}

	return cfg, nil
}

// Helper function to fetch a mandatory environment variable
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
