package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

type ExpirationChecker struct {
	metadataRepo  repository.MetadataRepository
	storageRepo   repository.StorageRepository
	kafkaProducer *kafka.KafkaProducer
}

func NewExpirationChecker(metadataRepo repository.MetadataRepository, storageRepo repository.StorageRepository, kafkaProducer *kafka.KafkaProducer) *ExpirationChecker {
	return &ExpirationChecker{
		metadataRepo:  metadataRepo,
		storageRepo:   storageRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (checker *ExpirationChecker) CheckForExpiredPastes(ctx context.Context, cfg *config.Config, log *jsonlog.Logger) {
	expiredKeys, err := checker.metadataRepo.DeleteAndReturnExpiredKeys()
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("Error retrieving expired pastes: %v", err), nil)
		return
	}

	if len(expiredKeys) == 0 {
		log.PrintInfo(ctx, "No expired pastes found.", nil)
		return
	}

	// Delete expired pastes from storage
	if err := checker.storageRepo.DeleteExpiredPastes(expiredKeys, cfg.BucketName); err != nil {
		log.PrintError(ctx, fmt.Errorf("Error deleting expired pastes from storage: %v", err), nil)
		return
	}

	// Format expired keys as JSON
	jsonExpiredKeys, err := formatExpiredKeysMessage(expiredKeys)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("Failed to encode keys to JSON: %v", err), nil)
		return
	}

	// Produce Kafka message for relocating expired keys
	if err := checker.kafkaProducer.ProduceMessages(jsonExpiredKeys, "relocate-key-topic"); err != nil {
		log.PrintError(ctx, fmt.Errorf("Failed to produce message to Kafka (Topic: %v): %v", "relocate-key-topic", err), nil)
		return
	}

	log.PrintInfo(ctx, fmt.Sprintf("Successfully deleted expired pastes and sent to Kafka: %v", expiredKeys), nil)
}

// Helper function to format expired keys as a JSON message
func formatExpiredKeysMessage(keys []string) (string, error) {
	message := struct {
		ExpiredKeys []string `json:"expired_keys"`
	}{
		ExpiredKeys: keys,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(msgBytes), nil
}
