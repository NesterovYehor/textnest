package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
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

	cfg.Addr = os.Getenv("PORT")
	cfg.RedisOption.Addr = os.Getenv("REDIS_ADDR")
	return cfg
}
