package services

import (
	"context"
	"fmt"

	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
)

type PasteService struct {
	metadataRepo repository.MetadataRepository
	storageRepo  repository.StorageRepository
}

func NewPasteService(metadataRepo repository.MetadataRepository, storageRepo repository.StorageRepository) *PasteService {
	return &PasteService{
		metadataRepo: metadataRepo,
		storageRepo:  storageRepo,
	}
}

func (service *PasteService) DeletePasteByKey(ctx context.Context, key string, bucketName string) error {
	if len(key) != 8 {
		return fmt.Errorf("key is not 8 characters: %s", key)
	}

	if err := service.metadataRepo.DeletePasteByKey(key); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	if err := service.storageRepo.DeletePasteByKey(key, bucketName); err != nil {
		return fmt.Errorf("failed to delete paste from storage: %w", err)
	}

	return nil
}
