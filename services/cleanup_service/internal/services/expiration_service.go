package services

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/scheduler"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/storage"
)

type ExpirationService struct {
	db *sql.DB
}

func NewExpirationService(db *sql.DB) *ExpirationService {
	return &ExpirationService{
		db: db,
	}
}

func (service *ExpirationService) Start(cfg *config.Config, ctx context.Context) {
	repo := repository.NewPasteRepository(service.db)
	storage, err := storage.NewS3Storage(cfg.BucketName, cfg.S3Region)
	if err != nil {
		log.Printf("Failed to connect to bucket %v, inn region %v", cfg.BucketName, cfg.S3Region)
		return
	}
	kafkaProducer, err := kafka.NewProducer(*cfg.Kafka, ctx)
	if err != nil {
		log.Println("Failed to create a kafka producer")
		return
	}

    defer kafkaProducer.Close()

	checker := scheduler.NewExpirationChecker(repo, storage, kafkaProducer)
	ticker := time.NewTicker(cfg.ExpirationInterval)

	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			checker.CheckForExpiredPastes(cfg)

		case <-stop:
			log.Println("Expiration service shutting down")
			kafkaProducer.Close()
			return
		}
	}
}
