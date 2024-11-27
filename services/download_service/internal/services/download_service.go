package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	logger "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/download_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DownloadService struct {
	metadataRepo  repository.MetadataRepository
	storageRepo   repository.StorageRepository
	metadataCache cache.Cache
	contentCache  cache.Cache
	config        *config.Config
	logger        *logger.Logger
	kafkaProducer *kafka.KafkaProducer
	pb.UnimplementedDownloadServiceServer
}

func NewDownloadService(
	storageRepo repository.StorageRepository,
	metadataRepo repository.MetadataRepository,
	logger *logger.Logger,
	config *config.Config,
	ctx context.Context,
) (*DownloadService, error) {
	if config.Kafka == nil || config.BucketName == "" {
		logger.PrintError(ctx, fmt.Errorf("invalid configuration"), nil)
		return nil, fmt.Errorf("invalid configuration: Kafka or BucketName missing")
	}

	metadataCache := cache.NewRedisCache(ctx, config.RedisMetadataAddr)
	contentCache := cache.NewRedisCache(ctx, config.RedisContentAddr)

	producer, err := kafka.NewProducer(*config.Kafka, ctx)
	if err != nil {
		logger.PrintError(ctx, err, nil)
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	return &DownloadService{
		metadataRepo:  metadataRepo,
		storageRepo:   storageRepo,
		logger:        logger,
		config:        config,
		kafkaProducer: producer,
		metadataCache: metadataCache,
		contentCache:  contentCache,
	}, nil
}

func (svc *DownloadService) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	// Validate the request key
	validator := validator.New()
	if svc.metadataRepo.IsKeyValid(req.Key, validator); !validator.Valid() {
		return nil, fmt.Errorf("invalid key")
	}

	// Fetch metadata
	metadata, err := svc.getMetadata(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Check expiration
	if err := svc.checkAndHandleExpiration(ctx, metadata); err != nil {
		return nil, err
	}

	// Fetch content
	content, err := svc.getContent(ctx, metadata.Key)
	if err != nil {
		return nil, err
	}

	// Build and return the response
	return svc.createDownloadResponse(req.Key, models.FromProto(metadata), content), nil
}

func (svc *DownloadService) getMetadata(ctx context.Context, key string) (*pb.Metadata, error) {
	// Try to fetch from cache
	cachedMetadata, found, err := svc.fetchMetadataFromCache(key)
	if err != nil {
		svc.logger.PrintError(ctx, err, map[string]string{"key": key})
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return cachedMetadata, nil
	}

	// Fetch from repository if not in cache
	dbMetadata, err := svc.metadataRepo.DownloadPasteMetadata(key)
	if err != nil {
		svc.logger.PrintError(ctx, err, map[string]string{"key": key})
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}
	rawData, err := proto.Marshal(dbMetadata.ToProto())
	if err != nil {
		return nil, err
	}

	// Cache the metadata
	if cacheErr := svc.metadataCache.Set(key, rawData, time.Minute*10); cacheErr != nil {
		svc.logger.PrintError(ctx, cacheErr, map[string]string{"key": key})
	}

	return dbMetadata.ToProto(), nil
}

func (svc *DownloadService) fetchMetadataFromCache(key string) (*pb.Metadata, bool, error) {
	serializedData, found, err := svc.metadataCache.Get(key)
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

func (svc *DownloadService) checkAndHandleExpiration(ctx context.Context, metadata *pb.Metadata) error {
	if isExpired(metadata.ExpiredDate.AsTime()) {
		err := svc.kafkaProducer.ProduceMessages(metadata.Key, "delete-expired-paste")
		if err != nil {
			svc.logger.PrintError(ctx, err, map[string]string{"key": metadata.Key})
			return fmt.Errorf("failed to produce delete-expired-paste message: %w", err)
		}

		svc.logger.PrintInfo(ctx, "Produced message to Kafka for expired paste", map[string]string{"key": metadata.Key})
		return fmt.Errorf("paste with key %s has expired", metadata.Key)
	}
	return nil
}

func (svc *DownloadService) getContent(ctx context.Context, key string) ([]byte, error) {
	// Try to fetch from cache
	content, found, err := svc.contentCache.Get(key)
	if err != nil {
		svc.logger.PrintError(ctx, err, map[string]string{"key": key})
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return content, nil
	}

	// Fetch from storage if not in cache
	content, err = svc.storageRepo.DownloadPasteContent(svc.config.BucketName, key)
	if err != nil {
		svc.logger.PrintError(ctx, err, map[string]string{"key": key})
		return nil, fmt.Errorf("could not download paste content: %w", err)
	}

	// Cache the content
	if cacheErr := svc.contentCache.Set(key, content, time.Hour); cacheErr != nil {
		svc.logger.PrintError(ctx, cacheErr, map[string]string{"key": key})
	}

	return content, nil
}

func (svc *DownloadService) createDownloadResponse(key string, metadata *models.Metadata, content []byte) *pb.DownloadResponse {
	return &pb.DownloadResponse{
		Key:            key,
		ExpirationDate: timestamppb.New(metadata.ExpiredDate),
		CreatedDate:    timestamppb.New(metadata.CreatedAt),
		Content:        content,
	}
}

func isExpired(expirationTime time.Time) bool {
	return time.Now().After(expirationTime)
}
