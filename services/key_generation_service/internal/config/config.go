package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
}

func InitConfig() *Config {
	cfg := &Config{} // Initialize cfg
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.Addr = os.Getenv("PORT")
	return cfg
}
