package config_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		logFile.Close()
		return
	}
	log := jsonlog.New(io.MultiWriter(logFile, os.Stdout), slog.LevelInfo)
	defer logFile.Close()

	ctx := context.Background()

	cfg, err := config.LoadConfig(log, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Grpc.Port)
	assert.Equal(t, "textnestbuycket", cfg.BucketName)
	assert.Equal(t, "eu-north-1", cfg.S3Region)
}
