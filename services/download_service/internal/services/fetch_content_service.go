package services

import (
	"context"
	"fmt"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
)

type FetchContentService struct {
	repo   *repository.ContentRepo
	logger *jsonlog.Logger
}

func NewFetchContentService(repo *repository.ContentRepo, log *jsonlog.Logger) (*FetchContentService, error) {
	return &FetchContentService{
		repo:   repo,
		logger: log,
	}, nil
}

func (svc *FetchContentService) GetContentUrl(ctx context.Context, key string) (string, error) {
	url, err := svc.repo.GenerateDownloadURL(key, ctx)
	if err != nil {
		svc.logger.PrintError(ctx, fmt.Errorf("Generattion down load url failed: %w", err), map[string]string{"key": key})
		return "", fmt.Errorf("could not generate url: %w", err)
	}

	return url, nil
}
