package config

import (
	"log"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc        *grpc.GrpcConfig
	RedisOption struct {
		Addr     string
		Password string
		DB       int
		Protocol int
	}
}

func InitConfig() *Config {
	cfg := &Config{} // Initialize cfg
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	cfg.RedisOption.Addr = os.Getenv("REDIS_ADDR")
	cfg.Grpc = &grpc.GrpcConfig{
		Port: port,
		Host: host,
	}
	return cfg
}
