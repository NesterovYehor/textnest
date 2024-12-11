package config_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../test_data/config.development.yaml")
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
    assert.Equal(t, ":8081", cfg.Grpc.Port)
	assert.Equal(t, "textnestbuycket", cfg.BucketName)
	assert.Equal(t, "eu-north-1", cfg.S3Region)
}
