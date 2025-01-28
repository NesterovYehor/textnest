package config

import (
	"context"
	"fmt"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/logger"
	"gopkg.in/yaml.v3"
)

// ServiceConfig holds gRPC connection details for all services.
type Config struct {
	UploadService   *grpc.GrpcConfig `yaml:"upload_service"`
	DownloadService *grpc.GrpcConfig `yaml:"download_service"`
	AuthService     *grpc.GrpcConfig `yaml:"auth_service"`
	KeyService      *grpc.GrpcConfig `yaml:"key_service"`
	HttpAddr        string           `yaml:"addr"`
}

// LoadConfig loads the gRPC service configuration from a YAML file.
func LoadConfig(ctx context.Context, log *jsonlog.Logger) (*Config, error) {
	// Retrieve CONFIG_PATH from environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}

	// Open the configuration file
	configFile, err := os.Open(configPath)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("unable to open configuration file: %w", err), nil)
		return nil, err
	}
	defer configFile.Close()

	// Decode YAML file into ServiceConfig
	var cfg Config
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to parse configuration file: %w", err), nil)
		return nil, err
	}

	// Validate required fields
	if cfg.UploadService == nil || cfg.UploadService.Port == "" {
		return nil, fmt.Errorf("upload service configuration is missing")
	}
	if cfg.KeyService == nil || cfg.KeyService.Port == "" {
		return nil, fmt.Errorf("key service configuration is missing")
	}

	return &cfg, nil
}
