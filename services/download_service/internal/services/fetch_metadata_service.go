package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/validation"
)

type FetchMetadataService struct {
	repo          *repository.MetadataRepo
	cache         cache.Cache
	kafkaProducer *kafka.KafkaProducer
}

func NewFetchMetadataService(repo *repository.MetadataRepo, cache cache.Cache, kafkaProducer *kafka.KafkaProducer) *FetchMetadataService {
	return &FetchMetadataService{
		repo:          repo,
		cache:         cache,
		kafkaProducer: kafkaProducer,
	}
}

func (svc *FetchMetadataService) FetchMetadataByKey(ctx context.Context, key string) (*pb.Metadata, error) {
	metadata, faund, err := svc.cache.Get(ctx, key)
	if faund && err == nil {
		return metadata, nil
	}
	if err != nil {
		return nil, err
	}
	metadata, err = svc.repo.DownloadPasteMetadata(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	if err := svc.checkAndHandleExpiration(ctx, metadata); err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := svc.cache.Set(ctx, key, metadata); err != nil {
			return
		}
	}()
	return metadata, nil
}

func (svc *FetchMetadataService) FetchMetadataByUserId(ctx context.Context, userId string, limit, offence int) ([]*pb.Metadata, error) {
	if err := validation.IsUserIdValid(userId); err != nil {
		return nil, err
	}

	metadataList, err := svc.repo.DownloadMetadataByUserId(ctx, userId, limit, offence)
	if err != nil {
		return nil, err
	}

	return metadataList, nil
}

func (svc *FetchMetadataService) checkAndHandleExpiration(ctx context.Context, metadata *pb.Metadata) error {
	if metadata.ExpiredDate.AsTime().After(time.Now()) {
		return nil
	}

	if err := svc.cache.Delete(ctx, metadata.Key); err != nil {
		return fmt.Errorf("failed to delete expired cache entry", err)
	}

	// Async Kafka message with retry
	go svc.retryKafkaMessage(metadata.Key, 3)

	return fmt.Errorf("paste with key '%s' has expired", metadata.Key)
}

func (svc *FetchMetadataService) retryKafkaMessage(key string, maxRetries int) {
	for i := 0; i < maxRetries; i++ {
		err := svc.kafkaProducer.ProduceMessages(key, "delete-expired-paste")
		if err == nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}
