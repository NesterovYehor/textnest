package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/scheduler"
)

type ExpirationService struct {
	db *sql.DB
}

func NewExpirationService(db *sql.DB) *ExpirationService {
	return &ExpirationService{
		db: db,
	}
}

func (service *ExpirationService) Start(cfg *config.Config, ctx context.Context, log *jsonlog.Logger) {
	// Initialize dependencies
	repo := repository.NewMetadataRepository(service.db)
	storage, err := repository.NewS3Storage(cfg.BucketName, cfg.S3Region)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("failed to create S3 storage: %w", err), nil)
		return
	}
	kafkaProducer, err := kafka.NewProducer(*cfg.Kafka, ctx)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("failed to create Kafka producer: %w", err), nil)
		return
	}

	checker := scheduler.NewExpirationChecker(repo, storage, kafkaProducer)
	ticker := time.NewTicker(cfg.ExpirationInterval)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer kafkaProducer.Close()
	defer ticker.Stop()
	log.PrintInfo(ctx, "Expiration Service started", nil)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			log.PrintInfo(ctx, "Starting expiration check", nil)
			checker.CheckForExpiredPastes(ctx, cfg, log)
			log.PrintInfo(ctx, "Expiration check completed", nil)

		case <-stopSignal:
			log.PrintInfo(ctx, "Received shutdown signal, cleaning up resources", nil)
			return
		}
	}
}
