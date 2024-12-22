package services

import (
	"context"
	"fmt"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/validation"
)

type MetadataManagementService struct {
	repo repository.MetadataRepository
	log  *jsonlog.Logger
}

func NewMetadataManagementService(repo repository.MetadataRepository, log *jsonlog.Logger) *MetadataManagementService {
	return &MetadataManagementService{repo: repo, log: log}
}

func (ms *MetadataManagementService) ValidateAndSave(ctx context.Context, metadata *models.MetaData) error {
	if v := validation.ValidateMetaData(metadata); !v.Valid() {
		err := fmt.Errorf("metadata validation errors: %v", v.Errors)
		ms.log.PrintError(ctx, err, map[string]string{"key": metadata.Key})
		return err
	}

	if err := ms.repo.UploadPasteMetadata(ctx, metadata); err != nil {
		err = fmt.Errorf("failed to save metadata: %w", err)
		ms.log.PrintError(ctx, err, map[string]string{"key": metadata.Key})
		return err
	}

	ms.log.PrintInfo(ctx, "Metadata saved successfully", map[string]string{"key": metadata.Key})
	return nil
}
