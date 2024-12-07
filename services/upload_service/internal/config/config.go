package config

import (
	"context"
	"fmt"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
)

const (
	DefaultGrpcPort = "4545"
	DefaultGrpcHost = "localhost"
	DefaultAppEnv   = "dev"
)

type Config struct {
	Grpc       *grpc.GrpcConfig
	BucketName string
	S3Region   string
	DBURL      string
}

// LoadConfig initializes the configuration by loading variables from the environment.
func LoadConfig(log *jsonlog.Logger, ctx context.Context) (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		log.PrintInfo(ctx, fmt.Sprintf("Grpc port not set, using default: %s", DefaultGrpcPort), nil)
		port = DefaultGrpcPort
	}
	host := os.Getenv("HOST")
	if host == "" {
		log.PrintInfo(ctx, fmt.Sprintf("Grpc host not set, using default: %s", DefaultGrpcHost), nil)
		host = DefaultGrpcHost
	}

	dbURL := getDatabaseURL(log, ctx)
	if dbURL == "" {
		return nil, fmt.Errorf("database URL not set")
	}

	bucketName := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("S3_REGION")
	if bucketName == "" || s3Region == "" {
		log.PrintError(ctx, fmt.Errorf("S3 configuration incomplete, some features may be unavailable"), nil)
	}

	return &Config{
		Grpc: &grpc.GrpcConfig{
			Port: port,
			Host: host,
		},
		BucketName: bucketName,
		S3Region:   s3Region,
		DBURL:      dbURL,
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
