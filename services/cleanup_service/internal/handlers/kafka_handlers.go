package handlers

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

func HandleDeleteExpiredPaste(message *sarama.ConsumerMessage, repo *repository.PasteRepository) error {
	key := string(message.Value)

	v := validator.New()

	if repo.IsKeyValid(v, key); !v.Valid() {
		return errors.New("Key is not 8 chars lenth")
	}

	err := repo.DeletePasteByKey(key)
	if err != nil {
		return err
	}

	return nil
}
