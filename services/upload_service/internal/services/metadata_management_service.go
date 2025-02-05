package services

import (
	"context"
	"fmt"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/validation"
)

type MetadataManagementService struct {
	repo *repository.MetadataRepository
	log  *jsonlog.Logger
}

func NewMetadataManagementService(repo *repository.MetadataRepository, log *jsonlog.Logger) *MetadataManagementService {
	return &MetadataManagementService{repo: repo, log: log}
}

func (ms *MetadataManagementService) ValidateAndSave(ctx context.Context, metadata *models.MetaData) error {
	if v := validation.ValidateMetaData(metadata); !v.Valid() {
		err := fmt.Errorf("metadata validation errors: %v", v.Errors)
		ms.log.PrintError(ctx, err, map[string]string{"key": metadata.Key})
		return err
	}

	if err := ms.repo.InsertPasteMetadata(ctx, metadata); err != nil {
		err = fmt.Errorf("failed to save metadata: %w", err)
		ms.log.PrintError(ctx, err, map[string]string{"key": metadata.Key})
		return err
	}

	ms.log.PrintInfo(ctx, "Metadata saved successfully", map[string]string{"key": metadata.Key})
	return nil
}

func (ms *MetadataManagementService) GetPasteOwner(ctx context.Context, key string) (string, error) {
	return ms.repo.GetPasteOwner(ctx, key)
}

func (ms *MetadataManagementService) UpdateMetadata(ctx context.Context, key string, expirationDate time.Time) error {
	return ms.repo.UpdatePasteMetadata(ctx, expirationDate, key)
}

func (ms *MetadataManagementService) ExpireMetadata(ctx context.Context, key string) error {
	return ms.repo.UpdatePasteMetadata(ctx, time.Now(), key)
}

func (ms *MetadataManagementService) ExpireAllPastes(ctx context.Context, userID string) error {
	return ms.repo.ExpireAllPastesByUserId(ctx, userID)
}
