package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
)

const (
	DefaultGrpcPort = "4545"
	DefaultGrpcHost = "localhost"
	DefaultAppEnv   = "dev"
)

type Config struct {
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	Kafka              *kafka.KafkaConfig
	DBUrl              string
}

// LoadConfig initializes the configuration by loading variables from the .env file and environment.
func LoadConfig(ctx context.Context, log *jsonlog.Logger) (*Config, error) {
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
	dbUrl := getDatabaseURL(log, ctx)
	if dbUrl == "" {
		log.PrintError(ctx, fmt.Errorf("DB_URL not set"), nil)
		return nil, fmt.Errorf("DB_URL not set")
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
		ExpirationInterval: interval,
		BucketName:         os.Getenv("S3_BUCKET_NAME"),
		S3Region:           os.Getenv("S3_REGION"),
		DBUrl:              dbUrl,
		Kafka: &kafka.KafkaConfig{
			Brokers:    brokers,
			Topics:     topics,
			MaxRetries: retries,
		},
	}, nil
}

func getDatabaseURL(log *jsonlog.Logger, ctx context.Context) string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = DefaultAppEnv
		log.PrintInfo(ctx, "APP_ENV not set, defaulting to 'dev'", nil)
	}

	var dbURL string
	switch env {
	case "dev":
		dbURL = os.Getenv("DB_URL_DEV")
	case "test":
		dbURL = os.Getenv("DB_URL_TEST")
	case "prod":
		dbURL = os.Getenv("DB_URL_PROD")
	default:
		log.PrintFatal(ctx, fmt.Errorf("Unknown APP_ENV: %s", env), nil)
		return ""
	}

	if dbURL == "" {
		log.PrintFatal(ctx, fmt.Errorf("Database URL not set for APP_ENV: %s", env), nil)
	}

	return dbURL
}
