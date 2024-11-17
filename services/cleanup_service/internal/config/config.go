package config

import (
	"os"
	"time"
)

type Config struct {
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	KafkaConfig        *KafkaConfig
}

func (cfg *Config) Init() {
	intervalStr := os.Getenv("EXPIRATION_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		interval = 5 * time.Minute
	}

	cfg.ExpirationInterval = interval
	cfg.BucketName = "textnestbuycket"
	cfg.S3Region = "eu-north-1"
    cfg.KafkaConfig = LoadKafkaConfig()
}
