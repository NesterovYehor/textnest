package download_service

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DownloadService struct {
	UnimplementedDownloadServiceServer
	dbRepo      repository.MetadataRepository
	storageRepo repository.StorageRepository
	cfg         *config.Config
	log         *jsonlog.Logger
	producer    *kafka.KafkaProducer
}

func NewDownloadService(
	storageRepo repository.StorageRepository,
	metadataRepo repository.MetadataRepository,
	log *jsonlog.Logger,
	cfg *config.Config,
	ctx context.Context,
) *DownloadService {
	producer, err := kafka.NewProducer(*cfg.Kafka, context.Background())
	if err != nil {
		log.PrintError(ctx, err, nil)
		return nil
	}
	return &DownloadService{
		dbRepo:      metadataRepo,
		storageRepo: storageRepo,
		log:         log,
		cfg:         cfg,
		producer:    producer,
	}
}

func (srv *DownloadService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error) {
	// Fetch metadata
	metadata, err := srv.dbRepo.DownloadPasteMetadata(req.Key)
	if err != nil {
		srv.log.PrintError(ctx, err, map[string]string{"key": req.Key})
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}

	// Handle expired pastes
	if isExpired(metadata.ExpiredDate) {
		err := srv.producer.ProduceMessages(metadata.Key, "delete-expired-paste")
		if err != nil {
			srv.log.PrintError(ctx, err, map[string]string{"key": req.Key})
			return nil, fmt.Errorf("failed to produce delete-expired-paste message: %w", err)
		}
		srv.log.PrintInfo(ctx, "Produced message to Kafka for expired paste", map[string]string{"key": req.Key})
		return nil, fmt.Errorf("paste with key %s has expired", req.Key)
	}

	// Fetch content
	content, err := srv.storageRepo.DownloadPasteContent(srv.cfg.BucketName, metadata.Key)
	if err != nil {
		srv.log.PrintError(ctx, err, map[string]string{"key": req.Key})
		return nil, fmt.Errorf("could not download paste content: %w", err)
	}

	// Return response
	return &DownloadResponse{
		Key:            req.Key,
		ExpirationDate: timestamppb.New(metadata.ExpiredDate),
		CreatedDate:    timestamppb.New(metadata.CreatedAt),
		Content:        content,
	}, nil
}

// Utility function for expiration check
func isExpired(expiredDate time.Time) bool {
	return time.Now().After(expiredDate)
}
