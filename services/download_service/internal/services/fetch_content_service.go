package services

import (
	"context"
	"fmt"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
)

type FetchContentService struct {
	repo       *repository.ContentRepo
	cache      cache.Cache
	logger     *jsonlog.Logger
	bucketName string // Optional if bucketName is constant for the service
}

func newFetchContentService(repo *repository.ContentRepo, cache cache.Cache, log *jsonlog.Logger, bucketName string) (*FetchContentService, error) {
	return &FetchContentService{
		repo:       repo,
		cache:      cache,
		logger:     log,
		bucketName: bucketName,
	}, nil
}

func (svc *FetchContentService) GetContent(ctx context.Context, key string) ([]byte, error) {
	// Try to fetch from cache
	content, found, err := svc.cache.Get(key)
	if err != nil {
		svc.logger.PrintError(ctx, fmt.Errorf("cache retrieval error: %w", err), map[string]string{"key": key})
	}

	if found {
		return content, nil
	}

	// Fetch from storage if not in cache
	content, err = svc.repo.DownloadPasteContent(svc.bucketName, key)
	if err != nil {
		svc.logger.PrintError(ctx, fmt.Errorf("storage retrieval error: %w", err), map[string]string{"key": key})
		return nil, fmt.Errorf("could not download paste content: %w", err)
	}

	// Cache the content asynchronously
	go func(key string, content []byte) {
		if cacheErr := svc.cache.Set(key, content, time.Hour); cacheErr != nil {
			svc.logger.PrintError(ctx, cacheErr, map[string]string{"key": key})
		}
	}(key, content)

	return content, nil
}
