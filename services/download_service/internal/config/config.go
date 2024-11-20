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
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc               *grpc.GrpcConfig
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	Kafka              *kafka.KafkaConfig
	DBURL              string
}

// LoadConfig initializes the configuration by loading variables from the .env file and environment.
func LoadConfig(log *jsonlog.Logger, ctx context.Context) (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.PrintError(ctx, err, nil)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.PrintError(ctx, fmt.Errorf("Grpc port not set, using default: 4444"), nil)
		port = "4444"
	}
	host := os.Getenv("HOST")
	if host == "" {
		log.PrintError(ctx, fmt.Errorf("Grpc host not set, using default: 4444"), nil)
		host = "localhost"
	}

	// Parse expiration interval
	intervalStr := os.Getenv("EXPIRATION_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil || interval == 0 {
		log.PrintError(ctx, fmt.Errorf("Invalid or missing EXPIRATION_INTERVAL. Using default: 5m. Error: %v", err), nil)
		interval = 5 * time.Minute
	}

	// Fetch Kafka brokers
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		log.PrintError(ctx, fmt.Errorf("KAFKA_BROKER not set, using default: localhost:9092"), nil)
		brokers[0] = "localhost:9092"
	}
	// Fetch Kafka brokers
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.PrintError(ctx, fmt.Errorf("DB_URL not set"), nil)
		return nil, fmt.Errorf("DB_URL not set")
	}

	// Fetch Kafka topics
	kafkaTopics := os.Getenv("KAFKA_TOPICS")
	var topics []string
	if kafkaTopics != "" {
		err := json.Unmarshal([]byte(kafkaTopics), &topics)
		if err != nil {
			log.PrintError(ctx, fmt.Errorf("failed to unmarshal KAFKA_TOPICS: %w", err), nil)
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
		log.PrintError(ctx, fmt.Errorf("Invalid or missing KAFKA_RETRIES. Using default: 5. Error: %v", err), nil)
		retries = 5
	}

	// Return the config struct with the parsed values
	return &Config{
		Grpc: &grpc.GrpcConfig{
			Port: port,
			Host: host,
		},
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
