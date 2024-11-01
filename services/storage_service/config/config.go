package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr    string
	Storage struct {
		Bucket    string
		Region    string
		AccessKey string // Capitalized for visibility
		SecretKey string // Changed to correct environment variable
	}
}

var (
	version   string
	buildTime string
)

func InitConfig() *Config {
	cfg := &Config{} // Initialize cfg
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.Addr = os.Getenv("PORT")
	cfg.Storage.Bucket = os.Getenv("BUCKET") // Correct key
	return cfg
}
