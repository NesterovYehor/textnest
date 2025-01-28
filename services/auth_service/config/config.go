package config

import (
	"errors"
	"os"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	SigningMethod jwt.SigningMethod
}

type Config struct {
	DBUrl     string
	Grpc      *grpc.GrpcConfig
	JwtConfig *JwtConfig
}

func LoadConfig(log *jsonlog.Logger) (*Config, error) {
	// Database URL
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, errors.New("database URL is not provided")
	}

	// gRPC Port
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		return nil, errors.New("gRPC port is not provided")
	}

	// JWT Secrets
	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		return nil, errors.New("access secret is not provided")
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, errors.New("refresh secret is not provided")
	}

	// JWT Expiries
	accessExpiry, err := parseDurationFromEnv("ACCESS_EXPIRY", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshExpiry, err := parseDurationFromEnv("REFRESH_EXPIRY", 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	signinMethod, err := getSigningMethodFromEnv()
	if err != nil {
		return nil, err
	}

	// Return the configuration
	return &Config{
		DBUrl: dbUrl,
		Grpc:  &grpc.GrpcConfig{Port: grpcPort},
		JwtConfig: &JwtConfig{
			AccessSecret:  accessSecret,
			RefreshSecret: refreshSecret,
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
			SigningMethod: signinMethod,
		},
	}, nil
}

// Helper function to parse durations from environment variables
func parseDurationFromEnv(envVar string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.New(envVar + " has an invalid duration format")
	}

	return duration, nil
}

func getSigningMethodFromEnv() (jwt.SigningMethod, error) {
	algorithm := os.Getenv("JWT_SIGNING_ALGORITHM")
	if algorithm == "" {
		return nil, errors.New("JWT_SIGNING_ALGORITHM is not set")
	}

	switch algorithm {
	case "HS256":
		return jwt.SigningMethodHS256, nil
	case "RS256":
		return jwt.SigningMethodRS256, nil
	default:
		return nil, errors.New("unsupported signing algorithm")
	}
}
