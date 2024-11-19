package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/pkg/validator"
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
	if err := checker.kafkaProducer.ProduceMessages(jsonExpiredKeys, cfg.Kafka.Topics["relocate-key-topic"]); err != nil {
		log.PrintError(ctx, fmt.Errorf("Failed to produce message to Kafka (Topic: %v): %v", cfg.Kafka.Topics["relocate-key-topic"], err), nil)
		return
	}

	log.PrintInfo(ctx, fmt.Sprintf("Successfully deleted expired pastes and sent to Kafka: %v", expiredKeys), nil)
}

func (checker *ExpirationChecker) DeletePasteByKey(message *sarama.ConsumerMessage, ctx context.Context, cfg *config.Config, log *jsonlog.Logger) error {
	key := string(message.Value)

	v := validator.New()

	// Validate that the key is 8 characters long
	if len(key) != 8 {
		err := errors.New(fmt.Sprintf("Key is not 8 chars length: %s", key))
		log.PrintError(ctx, err, nil)
		return err
	}

	// Validate key using the validator (if any additional validation is needed)
	if !v.Valid() {
		err := errors.New("Invalid key format")
		log.PrintError(ctx, err, nil)
		return err
	}

	// Attempt to delete the paste
	err := checker.metadataRepo.DeletePasteByKey(key)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return err
	}

	err = checker.storageRepo.DeletePasteByKey(key, cfg.BucketName)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return err
	}

	log.PrintInfo(ctx, "Expired paste is deleted", nil)

	return nil
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
