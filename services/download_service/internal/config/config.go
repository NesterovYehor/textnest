package config

import (
	"log"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc    *grpc.GrpcConfig
	Storage struct {
		Bucket string
		Region string
	}
	DbUrl               string
	KafkaConsumerConfig *KafkaConsumerConfig
}

func InitConfig() *Config {
	cfg := &Config{} // Initialize cfg
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	cfg.Storage.Bucket = os.Getenv("BUCKET")
	cfg.Storage.Region = os.Getenv("REGION")
	cfg.DbUrl = os.Getenv("DB_URL")
	cfg.Grpc = &grpc.GrpcConfig{
		Port: port,
		Host: host,
	}
	cfg.KafkaConsumerConfig = LoadKafkaConfig()
	return cfg
}
