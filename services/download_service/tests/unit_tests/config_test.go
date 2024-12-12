package unittests

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../test_data/config_test.yaml")
	defer os.Unsetenv("CONFIG_PATH")
	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	log := jsonlog.New(io.MultiWriter(logFile, os.Stdout), slog.LevelInfo)
	if log == nil {
		t.Fatalf("Logger initialization failed")
	}

	ctx := context.Background()

	cfg, err := config.LoadConfig(log, ctx)
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}
	assert.NotNil(t, cfg)

	assert.Equal(t, ":0000", cfg.Grpc.Port)
	assert.Equal(t, "localhost:0000", cfg.Kafka.Brokers[0])
	assert.Equal(t, "test-kafka-topic", cfg.Kafka.Topics[0])
	assert.Equal(t, "testbucket", cfg.BucketName)
	assert.Equal(t, "test-region", cfg.S3Region)
	assert.Equal(t, "test_db_url", cfg.DBURL)
	assert.Equal(t, ":0000", cfg.RedisMetadataAddr)
	assert.Equal(t, ":0000", cfg.RedisContentAddr)
	assert.Equal(t, 0, cfg.Kafka.MaxRetries)
	assert.Equal(t, time.Duration(0), cfg.ExpirationInterval)
}
