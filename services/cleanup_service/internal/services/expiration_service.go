package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

type ExpirationService struct {
	metadataRepo  repository.MetadataRepository
	storageRepo   repository.StorageRepository
	kafkaProducer *kafka.KafkaProducer
	bucketName    string
}

func NewExpirationService(
	metadataRepo repository.MetadataRepository,
	storageRepo repository.StorageRepository,
	kafkaProducer *kafka.KafkaProducer,
	bucketName string,
) *ExpirationService {
	return &ExpirationService{
		metadataRepo:  metadataRepo,
		storageRepo:   storageRepo,
		kafkaProducer: kafkaProducer,
		bucketName:    bucketName,
	}
}

func (s *ExpirationService) ProcessExpirations(ctx context.Context) error {
	// Step 1: Retrieve expired keys
	expiredKeys, err := s.metadataRepo.DeleteAndReturnExpiredKeys()
	if err != nil {
		return fmt.Errorf("error retrieving expired pastes: %v", err)
	}

	if len(expiredKeys) == 0 {
		return nil // No expired keys to process
	}

	// Step 2: Delete expired keys from storage
	if err := s.storageRepo.DeleteExpiredPastes(expiredKeys, s.bucketName); err != nil {
		return fmt.Errorf("error deleting expired pastes from storage: %v", err)
	}

	// Step 3: Send expired keys to Kafka
	jsonExpiredKeys, err := json.Marshal(map[string][]string{"expired_keys": expiredKeys})
	if err != nil {
		return fmt.Errorf("failed to encode keys to JSON: %v", err)
	}

	if err := s.kafkaProducer.ProduceMessages(string(jsonExpiredKeys), "relocate-key-topic"); err != nil {
		return fmt.Errorf("failed to produce message to Kafka: %v", err)
	}

	return nil
}
