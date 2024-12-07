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
)

const (
	DefaultGrpcPort = "4545"
	DefaultGrpcHost = "localhost"
	DefaultAppEnv   = "dev"
)

type Config struct {
	Grpc               *grpc.GrpcConfig
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	Kafka              *kafka.KafkaConfig
	DBURL              string
	RedisMetadataAddr  string
	RedisContentAddr   string
}

// LoadConfig loads configuration values from environment variables and the .env file.
func LoadConfig(log *jsonlog.Logger, ctx context.Context) (*Config, error) {
	// Fetch and log GRPC server configuration (host and port)
	grpcHost, grpcPort := getGRPCConfig(log, ctx)

	// Fetch and parse expiration interval (default 5 minutes if invalid)
	expirationInterval := getExpirationInterval(log, ctx)

	// Fetch and validate Kafka brokers and retry settings
	kafkaConfig := getKafkaConfig(log, ctx)

	// Fetch and validate database URL (must be provided)
	dbURL := getDatabaseURL(log, ctx)

	// Fetch and validate Redis addresses for metadata and content caches
	redisMetadataAddr := os.Getenv("METADATA_CACHE_REDIS_HOST")
	redisContentAddr := os.Getenv("CONTENT_CACHE_REDIS_HOST")

	// Fetch and validate S3 bucket name and region
	bucketName := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("S3_REGION")

	// Return populated config struct
	return &Config{
		Grpc: &grpc.GrpcConfig{
			Port: grpcPort,
			Host: grpcHost,
		},
		RedisMetadataAddr:  redisMetadataAddr,
		RedisContentAddr:   redisContentAddr,
		ExpirationInterval: expirationInterval,
		DBURL:              dbURL,
		BucketName:         bucketName,
		S3Region:           s3Region,
		Kafka:              kafkaConfig,
	}, nil
}

// getGRPCConfig fetches GRPC host and port, with default values if missing.
func getGRPCConfig(log *jsonlog.Logger, ctx context.Context) (string, string) {
	host := os.Getenv("HOST")
	if host == "" {
		log.PrintError(ctx, fmt.Errorf("GRPC host not set, using default: localhost"), nil)
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.PrintError(ctx, fmt.Errorf("GRPC port not set, using default: 4444"), nil)
		port = "4444"
	}
	return host, port
}

// getExpirationInterval fetches the expiration interval, defaulting to 5 minutes if invalid.
func getExpirationInterval(log *jsonlog.Logger, ctx context.Context) time.Duration {
	intervalStr := os.Getenv("EXPIRATION_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil || interval == 0 {
		log.PrintError(ctx, fmt.Errorf("Invalid or missing EXPIRATION_INTERVAL. Using default: 5m. Error: %v", err), nil)
		interval = 5 * time.Minute
	}
	return interval
}

// getKafkaConfig fetches Kafka configuration including brokers, topics, and retries.
func getKafkaConfig(log *jsonlog.Logger, ctx context.Context) *kafka.KafkaConfig {
	// Fetch Kafka brokers
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		log.PrintError(ctx, fmt.Errorf("KAFKA_BROKER not set, using default: localhost:9092"), nil)
		brokers[0] = "localhost:9092"
	}

	// Parse Kafka topics
	kafkaTopics := os.Getenv("KAFKA_TOPICS")
	var topics []string
	if kafkaTopics != "" {
		err := json.Unmarshal([]byte(kafkaTopics), &topics)
		if err != nil {
			log.PrintError(ctx, fmt.Errorf("Failed to unmarshal KAFKA_TOPICS: %w", err), nil)
			topics = []string{}
		}
	}

	// Fetch Kafka retry settings
	retryStr := os.Getenv("KAFKA_RETRIES")
	retries, err := strconv.Atoi(retryStr)
	if err != nil || retries <= 0 {
		log.PrintError(ctx, fmt.Errorf("Invalid or missing KAFKA_RETRIES. Using default: 5. Error: %v", err), nil)
		retries = 5
	}

	return &kafka.KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		MaxRetries: retries,
	}
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
