package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/joho/godotenv"
)

type Config struct {
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	Kafka              *kafka.KafkaConfig
}

// LoadConfig initializes the configuration by loading variables from the .env file and environment.
func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse expiration interval
	intervalStr := os.Getenv("EXPIRATION_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil || interval == 0 {
		log.Printf("Invalid or missing EXPIRATION_INTERVAL. Using default: 5m. Error: %v", err)
		interval = 5 * time.Minute
	}

	// Fetch Kafka brokers
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		log.Println("KAFKA_BROKER not set, using default: localhost:9092")
		brokers[0] = "localhost:9092"
	}

	// Fetch Kafka topics
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

	// Kafka retry settings
	retryStr := os.Getenv("KAFKA_RETRIES")
	retries, err := strconv.Atoi(retryStr)
	if err != nil || retries <= 0 {
		log.Printf("Invalid or missing KAFKA_RETRIES. Using default: 5. Error: %v", err)
		retries = 5
	}

	// Return the config struct with the parsed values
	return &Config{
		ExpirationInterval: interval,
		BucketName:         os.Getenv("S3_BUCKET_NAME"),
		S3Region:           os.Getenv("S3_REGION"),
		Kafka: &kafka.KafkaConfig{
			Brokers:    brokers,
			Topics:     topics,
			MaxRetries: retries,
		},
	}, nil
}
