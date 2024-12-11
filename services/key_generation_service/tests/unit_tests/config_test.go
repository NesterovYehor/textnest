package config_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
)

func TestLoadConfig(t *testing.T) {
	// Set up logging
	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log := jsonlog.New(multiWriter, slog.LevelInfo)
	assert.NotNil(t, log)

	// Load configuration
	cfg, err := config.LoadConfig(context.Background(), log)
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	// Validate configuration values
	assert.Equal(t, ":5555", cfg.Grpc.Port)
	assert.Equal(t, "localhost:6378", cfg.RedisOption.Addr)
	assert.Equal(t, "localhost:9092", cfg.Kafka.Brokers[0])
	assert.Equal(t, 0, cfg.Kafka.MaxRetries)
	assert.Contains(t, cfg.Kafka.Topics, "delete-expired-key-topic")
	assert.Contains(t, cfg.Kafka.Topics, "user-notifications")
	assert.Equal(t, "5m0s", cfg.ExpirationInterval.String())
}
