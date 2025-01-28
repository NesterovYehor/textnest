package config

import (
	"context"
	"fmt"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "config.yaml"

type Config struct {
	Grpc       *grpc.GrpcConfig `yaml:"grpc"`
	BucketName string           `yaml:"bucket_name"`
	S3Region   string           `yaml:"region"`
	DBURL      string           `yaml:"db"`
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(log *jsonlog.Logger, ctx context.Context) (*Config, error) {
	// Read CONFIG_PATH from environment
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yaml" // Default path
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
	if cfg.DBURL == "" {
		log.PrintFatal(ctx, fmt.Errorf("database URL is not set"), nil)
	}
	if cfg.BucketName == "" || cfg.S3Region == "" {
		log.PrintError(ctx, fmt.Errorf("S3 configuration is incomplete, some features may be unavailable"), nil)
	}

	return &cfg, nil
}
