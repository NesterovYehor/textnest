package services

import (
	"context"
	"fmt"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
)

type ContentManagementService struct {
	repo *repository.ContentRepository
	log  *jsonlog.Logger
}

func NewStorageService(repo *repository.ContentRepository, log *jsonlog.Logger) *ContentManagementService {
	return &ContentManagementService{repo: repo, log: log}
}

func (svc *ContentManagementService) SaveContent(ctx context.Context, key string, data []byte) error {
	if err := svc.repo.UploadPasteContent(ctx, key, data); err != nil {
		err = fmt.Errorf("failed to save content: %w", err)
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
		return err
	}

	svc.log.PrintInfo(ctx, "Content uploaded successfully", map[string]string{"key": key})
	return nil
}

func (svc *ContentManagementService) UpdateContent(ctx context.Context, key string, data []byte) error {
	if err := svc.repo.UploadPasteContent(ctx, key, data); err != nil {
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
		return err
	}
	return nil
}
