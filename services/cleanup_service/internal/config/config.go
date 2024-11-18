package config

import (
	"os"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
)

type Config struct {
	ExpirationInterval time.Duration
	BucketName         string
	S3Region           string
	Kafka              *kafka.KafkaConfig
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
	brokers := make([]string, 0, 1)
	brokers = append(brokers, "localhost:9092")
	topics := make([]string, 0, 1)
	topics = append(brokers, "relocate-key")
	cfg.Kafka = kafka.LoadKafkaConfig(brokers, topics, "1", 5)
}
