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

func (svc *ContentManagementService) GenerateUploadURL(ctx context.Context, key string) (string, error) {
	// Validate input
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	// Call the repository to get a presigned URL
	uploadURL, err := svc.repo.GenerateUploadURL(ctx, key)
	if err != nil {
		err = fmt.Errorf("failed to generate upload URL: %w", err)
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
		return "", err
	}

	// Log success
	svc.log.PrintInfo(ctx, "Presigned URL generated successfully", map[string]string{"key": key})

	return uploadURL, nil
}
