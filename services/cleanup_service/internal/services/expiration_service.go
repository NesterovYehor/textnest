package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

type ExpirationService struct {
	metadataRepo  *repository.MetadataRepo
	storageRepo   *repository.StorageRepo
	kafkaProducer *kafka.KafkaProducer
}

func NewExpirationService(
	metadataRepo *repository.MetadataRepo,
	storageRepo *repository.StorageRepo,
	kafkaProducer *kafka.KafkaProducer,
) *ExpirationService {
	return &ExpirationService{
		metadataRepo:  metadataRepo,
		storageRepo:   storageRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *ExpirationService) ProcessExpirations(ctx context.Context) error {
	// Step 1: Retrieve expired keys
	expiredKeys, err := s.metadataRepo.DeleteAndReturnExpiredKeys()
	if err != nil {
		return fmt.Errorf("error retrieving expired pastes: %v", err)
	}
	log.Println(expiredKeys)

	if len(expiredKeys) == 0 {
		return nil // No expired keys to process
	}

	// Step 2: Delete expired keys from storage
	if err := s.storageRepo.DeleteExpiredPastes(expiredKeys); err != nil {
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
