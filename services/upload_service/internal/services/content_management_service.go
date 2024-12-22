package services

import (
	"context"
	"fmt"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
)

type ContentManagementService struct {
	repo repository.StorageRepository
	log  *jsonlog.Logger
}

func NewStorageService(repo repository.StorageRepository, log *jsonlog.Logger) *ContentManagementService {
	return &ContentManagementService{repo: repo, log: log}
}

func (ss *ContentManagementService) SaveContent(ctx context.Context, bucket, key string, data []byte) error {
	if err := ss.repo.UploadPasteContent(ctx, bucket, key, data); err != nil {
		err = fmt.Errorf("failed to save content: %w", err)
		ss.log.PrintError(ctx, err, map[string]string{"key": key})
		return err
	}

	ss.log.PrintInfo(ctx, "Content uploaded successfully", map[string]string{"key": key})
	return nil
}
