package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/gogo/protobuf/proto"
)

type FetchMetadataService struct {
	repo          repository.MetadataRepository
	cache         cache.Cache
	kafkaProducer *kafka.KafkaProducer
	log           *jsonlog.Logger
}

func newFetchMetadataService(repo repository.MetadataRepository, cache cache.Cache, log *jsonlog.Logger, kafkaProducer *kafka.KafkaProducer) (*FetchMetadataService, error) {
	return &FetchMetadataService{
		repo:          repo,
		cache:         cache,
		kafkaProducer: kafkaProducer,
		log:           log,
	}, nil
}

func (svc *FetchMetadataService) GetMetadata(ctx context.Context, key string) (*pb.Metadata, error) {
	// Try to fetch from cache
	cachedMetadata, found, err := svc.fetchMetadataFromCache(key)
	if err != nil {
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
	}
	if found {
		return cachedMetadata, nil
	}

	// Fetch from repository if not in cache
	dbMetadata, err := svc.repo.DownloadPasteMetadata(key)
	if err != nil {
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}
	rawData, err := proto.Marshal(dbMetadata.ToProto())
	if err != nil {
		return nil, err
	}

	// Cache the metadata
	if cacheErr := svc.cache.Set(key, rawData, time.Minute*10); cacheErr != nil {
		svc.log.PrintError(ctx, cacheErr, map[string]string{"key": key})
	}

	return dbMetadata.ToProto(), nil
}

func (svc *FetchMetadataService) fetchMetadataFromCache(key string) (*pb.Metadata, bool, error) {
	serializedData, found, err := svc.cache.Get(key)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	var metadata pb.Metadata
	if err := proto.Unmarshal(serializedData, &metadata); err != nil {
		return nil, true, err
	}
	return &metadata, true, nil
}

func (svc *FetchMetadataService) checkAndHandleExpiration(ctx context.Context, metadata *pb.Metadata) error {
	if metadata.ExpiredDate.AsTime().Before(time.Now()) {
		err := svc.kafkaProducer.ProduceMessages(metadata.Key, "delete-expired-paste")
		if err != nil {
			svc.log.PrintError(ctx, err, map[string]string{"key": metadata.Key})
			return fmt.Errorf("failed to produce delete-expired-paste message: %w", err)
		}

		svc.log.PrintInfo(ctx, "Produced message to Kafka for expired paste", map[string]string{"key": metadata.Key})
		return fmt.Errorf("paste has expired")
	}
	return nil
}
