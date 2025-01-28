package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ExpirationInterval time.Duration      `yaml:"expiration_interval"`
	BucketName         string             `yaml:"bucket_name"`
	Kafka              *kafka.KafkaConfig `yaml:"kafka"`
	DBUrl              string             `yaml:"db_url"`
}

// LoadConfig initializes the configuration by loading variables from the .env file and environment.
func LoadConfig(ctx context.Context) (*Config, error) {
	// Read CONFIG_PATH from environment
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}
	data, err := os.Open(path)
	if err != nil {
        return nil, fmt.Errorf("failed to read configuration file: %w, on path: %v", err, path)
	}
	defer data.Close()

	var cfg Config
	decoder := yaml.NewDecoder(data)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}
	// Validate required fields
	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("gRPC configuration is incomplete")
	}
	if cfg.ExpirationInterval < time.Second || cfg.ExpirationInterval > time.Hour {
		return nil, fmt.Errorf("Timeout duration should be between 1 second and 1 hour, got: %v", cfg.ExpirationInterval)
	}
	if cfg.Kafka == nil || len(cfg.Kafka.Topics) == 0 || len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("kafka configuration is incomplete")
	}
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("S3 configuration is incomplete")
	}

	return &cfg, nil
}
