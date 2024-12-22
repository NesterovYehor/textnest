package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Grpc              *grpc.GrpcConfig  `yaml:"grpc"`
	Kafka             kafka.KafkaConfig `yaml:"kafka"`
	BucketName        string            `yaml:"bucket_name"`
	S3Region          string            `yaml:"region"`
	DBURL             string            `yaml:"db_url"`
	RedisMetadataAddr string            `yaml:"metadata_cache_addr"`
	RedisContentAddr  string            `yaml:"content_cache_addr"`

	ExpirationInterval time.Duration `yaml:"expiration_interval"`
}

// LoadConfig loads configuration values from environment variables and the .env file.
func LoadConfig(log *jsonlog.Logger, ctx context.Context) (*Config, error) {
	// Read CONFIG_PATH from environment
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}
	data, err := os.Open(path)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to read configuration file: %w", err), nil)
		return nil, err
	}
	defer data.Close()

	var cfg Config
	decoder := yaml.NewDecoder(data)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to parse configuration file: %w", err), nil)
		return nil, err
	}
	// Validate required fields
	if cfg.Grpc == nil || cfg.Grpc.Port == "" {
		log.PrintFatal(ctx, fmt.Errorf("gRPC configuration is incomplete"), nil)
	}
	if len(cfg.Kafka.Topics) == 0 || len(cfg.Kafka.Brokers) == 0 {
		log.PrintFatal(ctx, fmt.Errorf("kafka configuration is incomplete"), nil)
	}
	if cfg.BucketName == "" || cfg.S3Region == "" {
		log.PrintInfo(ctx, cfg.BucketName, nil)
		log.PrintInfo(ctx, cfg.S3Region, nil)
		log.PrintFatal(ctx, fmt.Errorf("S3 configuration is incomplete"), nil)
	}
	if cfg.RedisContentAddr == "" || cfg.RedisMetadataAddr == "" {
		log.PrintFatal(ctx, fmt.Errorf("redis cahce configuration is incomplete"), nil)
	}
	return &cfg, nil
}
