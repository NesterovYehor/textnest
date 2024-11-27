package config

import (
	"context"
	"fmt"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc       *grpc.GrpcConfig
	BucketName string
	S3Region   string
	DBURL      string
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

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.PrintError(ctx, fmt.Errorf("DB_URL not set"), nil)
		return nil, fmt.Errorf("DB_URL not set")
	}

	// Return the config struct with the parsed values
	return &Config{
		Grpc: &grpc.GrpcConfig{
			Port: port,
			Host: host,
		},
		BucketName: os.Getenv("S3_BUCKET_NAME"),
		S3Region:   os.Getenv("S3_REGION"),
		DBURL:      dbUrl,
	}, nil
}
