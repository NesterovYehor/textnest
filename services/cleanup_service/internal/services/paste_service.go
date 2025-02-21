package services

import (
	"context"
	"fmt"

	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

type PasteService struct {
	metadataRepo *repository.MetadataRepo
	storageRepo  *repository.StorageRepo
}

func NewPasteService(metadataRepo *repository.MetadataRepo, storageRepo *repository.StorageRepo) *PasteService {
	return &PasteService{
		metadataRepo: metadataRepo,
		storageRepo:  storageRepo,
	}
}

func (service *PasteService) DeletePasteByKey(ctx context.Context, key string) error {
	if len(key) != 8 {
		return fmt.Errorf("key is not 8 characters: %s", key)
	}

	if err := service.metadataRepo.DeletePasteByKey(key); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	if err := service.storageRepo.DeletePasteContentByKey(key); err != nil {
		return fmt.Errorf("failed to delete paste from storage: %w", err)
	}

	return nil
}
