package scheduler

import (
	"encoding/json"
	"log"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/storage"
)

type ExpirationChecker struct {
	repo          *repository.PasteRepository
	storage       storage.Storage // Use the interface directly, not a pointer to it
	kafkaProducer *kafka.KafkaProducer
}

func NewExpirationChecker(repo *repository.PasteRepository, storage storage.Storage, kafkaProducer *kafka.KafkaProducer) *ExpirationChecker {
	return &ExpirationChecker{
		repo:          repo,
		storage:       storage,
		kafkaProducer: kafkaProducer,
	}
}

func (checker *ExpirationChecker) CheckForExpiredPastes(cfg *config.Config) {
	expiredKeys, err := checker.repo.DeleteAndReturnExpiredKeys()
	if err != nil {
		log.Printf("Error retrieving expired pastes: %v", err)
		return
	}

	if len(expiredKeys) == 0 {
		log.Println("No expired pastes found.")
		return
	}

	if err := checker.storage.DeleteExpiredPastes(expiredKeys); err != nil {
		log.Printf("Error deleting expired pastes from storage: %v", err)
		return
	}

	jsonExpiredKeys, err := formatExpiredKeysMessage(expiredKeys)
	if err != nil {
		log.Printf("Failed to encode keys to JSON: %v", err)
		return
	}

	if err := checker.kafkaProducer.ProduceMessages(jsonExpiredKeys, "Relocate-Keys-Topic"); err != nil {
		log.Printf("Failed to produce message to Kafka (Topic: %s): %v", "Relocate-Keys-Topic", err)
		return
	}
	log.Printf("Successfully deleted expired pastes and sent to Kafka: %v", expiredKeys)
}

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

	return string(msgBytes), err
}
