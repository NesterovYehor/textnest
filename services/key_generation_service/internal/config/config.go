package config

import (
	"log"
	"os"

	httpserver "github.com/NesterovYehor/TextNest/pkg/http"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc        *httpserver.Config
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
	cfg.RedisOption.Addr = os.Getenv("REDIS_ADDR")
	cfg.Grpc = grpc
	return cfg
}
