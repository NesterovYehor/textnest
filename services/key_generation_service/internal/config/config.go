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
	Grpc      *grpc.GrpcConfig  `yaml:"grpc"`
	Kafka     kafka.KafkaConfig `yaml:"kafka"`
	RedisAddr string            `yaml:"redis"`

	ExpirationInterval time.Duration `yaml:"expirationInterval"`
}

func LoadConfig(ctx context.Context, log *jsonlog.Logger) (*Config, error) {
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
		log.PrintInfo(ctx, fmt.Sprintf("%+v\n", cfg.Kafka), nil)
		log.PrintFatal(ctx, fmt.Errorf("kafka configuration is incomplete"), nil)
	}
	if cfg.RedisAddr == "" {
		log.PrintFatal(ctx, fmt.Errorf("redis configuration is incomplete"), nil)
	}

	return &cfg, nil
}
