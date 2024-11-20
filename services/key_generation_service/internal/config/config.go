package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
)

type Config struct {
	Grpc        *grpc.GrpcConfig
	Kafka       *kafka.KafkaConfig
	RedisOption struct {
		Addr     string
		Password string
		DB       int
	}
	ExpirationInterval time.Duration
}

func InitConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}
	cfg := &Config{}

	// Redis Config
	cfg.RedisOption.Addr = getEnvOrFatal("REDIS_ADDR")
	cfg.RedisOption.Password = os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	if dbStr != "" {
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
		}
		cfg.RedisOption.DB = db
	}

	// gRPC Config
	port := getEnvOrFatal("PORT")
	host := getEnvOrFatal("HOST")
	cfg.Grpc = &grpc.GrpcConfig{
		Port: port,
		Host: host,
	}

	// Expiration Interval
	intervalStr := os.Getenv("EXPIRATION_INTERVAL")
	if intervalStr == "" {
		return nil, fmt.Errorf("EXPIRATION_INTERVAL is required")
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid EXPIRATION_INTERVAL: %w", err)
	}
	cfg.ExpirationInterval = interval

	// Kafka Config
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		log.Println("KAFKA_BROKER not set, using default: localhost:9092")
		brokers[0] = "localhost:9092"
	}
	kafkaTopics := os.Getenv("KAFKA_TOPICS")
	var topics []string
	if kafkaTopics != "" {
		err := json.Unmarshal([]byte(kafkaTopics), &topics)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal KAFKA_TOPICS: %w", err)
		}
	} else {
		// Default to empty map if no topics are specified
		topics = make([]string, 0, 1)
	}
	retryStr := os.Getenv("KAFKA_RETRIES")
	retries, err := strconv.Atoi(retryStr)
	if err != nil || retries <= 0 {
		log.Printf("Invalid or missing KAFKA_RETRIES. Using default: 5. Error: %v", err)
		retries = 5
	}
	cfg.Kafka = &kafka.KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		MaxRetries: retries,
	}

	return cfg, nil
}

// Helper function to fetch a mandatory environment variable
func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}
