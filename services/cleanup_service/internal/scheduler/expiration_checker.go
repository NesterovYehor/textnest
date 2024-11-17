package scheduler

import (
	"encoding/json"
	"log"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
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

func (checker *ExpirationChecker) CheckForExpiredPastes() {
	// Retrieve expired keys from the repository
	expiredKeys, err := checker.repo.DeleteAndReturnExpiredKeys()
	if err != nil {
		log.Println("Error retrieving expired pastes:", err)
		return
	}

	// Delete expired pastes from storage
	err = checker.storage.DeleteExpiredPastes(expiredKeys) // Pass appropriate lifetimeSecs
	if err != nil {
		log.Println("Error deleting expired pastes from storage:", err)
		return
	}

	jsonExpiredKeys, err := formatExpiredKeysMessage(expiredKeys)
	if err != nil {
		log.Println("Failded to encode keys to json")
		return
	}

	err = kafka.
	if err != nil {
		log.Println("Failed to produce message to kafka")
		return
	}

	// Log the deleted keys
	if len(expiredKeys) > 0 {
		log.Printf("Successfully deleted expired pastes: %v", expiredKeys)
	} else {
		log.Println("No expired pastes found.")
	}
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
