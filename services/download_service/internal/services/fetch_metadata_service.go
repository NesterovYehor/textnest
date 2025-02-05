package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/validation"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FetchMetadataService struct {
	repo          *repository.MetadataRepo
	cache         cache.Cache
	kafkaProducer *kafka.KafkaProducer
	log           *jsonlog.Logger
}

func newFetchMetadataService(repo *repository.MetadataRepo, cache cache.Cache, log *jsonlog.Logger, kafkaProducer *kafka.KafkaProducer) (*FetchMetadataService, error) {
	return &FetchMetadataService{
		repo:          repo,
		cache:         cache,
		kafkaProducer: kafkaProducer,
		log:           log,
	}, nil
}

func (svc *FetchMetadataService) FetchMetadataByKey(ctx context.Context, key string) (*pb.Metadata, error) {
	// Try cache first
	cachedMetadata, found, err := svc.fetchMetadataFromCache(key)
	if err != nil {
		svc.log.PrintError(ctx, err, map[string]string{"key": key})
	}
	if found {
		// Check expiration even for cached entries
		if err := svc.checkAndHandleExpiration(ctx, cachedMetadata); err != nil {
			return nil, err
		}
		return cachedMetadata, nil
	}

	// Fetch from repository
	dbMetadata, err := svc.repo.DownloadPasteMetadata(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	// Check expiration
	pbMetadata := dbMetadata.ToProto()
	if err := svc.checkAndHandleExpiration(ctx, pbMetadata); err != nil {
		return nil, err
	}

	// Update cache (async to avoid blocking)
	go func() {
		rawData, err := proto.Marshal(pbMetadata)
		if err != nil {
			svc.log.PrintError(ctx, err, map[string]string{"key": key})
			return
		}
		if err := svc.cache.Set(key, rawData, 10*time.Minute); err != nil {
			svc.log.PrintError(ctx, err, map[string]string{"key": key})
		}
	}()

	return pbMetadata, nil
}

func (svc *FetchMetadataService) FetchMetadataByUserId(ctx context.Context, userId string, limit, offence int) ([]*pb.Metadata, error) {
	if err := validation.IsUserIdValid(userId); err != nil {
		return nil, err
	}

	// Fetch metadata from the repository
	metadataList, err := svc.repo.DownloadMetadataByUserId(ctx, userId, limit, offence)
	if err != nil {
		return nil, err
	}

	// Convert to gRPC Metadata format
	grpcMetadataList := ConvertMetadataListToProto(metadataList)

	// Return gRPC response without totalCount
	return grpcMetadataList, nil
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
	if metadata.ExpiredDate.AsTime().After(time.Now()) {
		return nil
	}

	// Invalidate cache
	if err := svc.cache.Delete(metadata.Key); err != nil {
		svc.log.PrintInfo(ctx, "failed to delete expired cache entry", map[string]string{
			"key": metadata.Key,
			"err": err.Error(),
		})
	}

	// Async Kafka message with retry
	go svc.retryKafkaMessage(ctx, metadata.Key, 3)

	return fmt.Errorf("paste with key '%s' has expired", metadata.Key)
}

func (svc *FetchMetadataService) retryKafkaMessage(ctx context.Context, key string, maxRetries int) {
	for i := 0; i < maxRetries; i++ {
		err := svc.kafkaProducer.ProduceMessages(key, "delete-expired-paste")
		if err == nil {
			svc.log.PrintInfo(ctx, "expired paste deletion queued", map[string]string{"key": key})
			return
		}
		svc.log.PrintError(ctx, err, map[string]string{
			"key":   key,
			"retry": fmt.Sprintf("%d/%d", i+1, maxRetries),
		})
		time.Sleep(2 * time.Second)
	}
}

func ConvertMetadataListToProto(metadataList []models.Metadata) []*pb.Metadata {
	var grpcMetadataList []*pb.Metadata

	for _, meta := range metadataList {
		grpcMeta := &pb.Metadata{
			Key: meta.Key,
		}
		if !meta.CreatedAt.IsZero() {
			grpcMeta.CreatedAt = timestamppb.New(meta.CreatedAt)
		}
		if !meta.ExpirationDate.IsZero() {
			grpcMeta.ExpiredDate = timestamppb.New(meta.ExpirationDate)
		}

		grpcMetadataList = append(grpcMetadataList, grpcMeta)
	}

	return grpcMetadataList
}
