package coordinators

import (
	"context"
	"database/sql"
	"sync"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	log "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/cache"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
)

type DownloadCoordinator struct {
	fetchMetadataService *services.FetchMetadataService
	fetchContentService  *services.FetchContentService
	cfg                  *config.Config
	logger               *log.Logger
	pb.UnsafePasteDownloadServer
}

func NewDownloadCoordinator(ctx context.Context, cfg *config.Config, log *log.Logger, db *sql.DB) (*DownloadCoordinator, error) {
	cache, err := cache.NewRedisCache(cfg.RedisMetadataAddr)
	if err != nil {
		return nil, err
	}
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka, ctx)
	if err != nil {
		return nil, err
	}
	metadataRepo := repository.NewMetadataRepo(db)
	contentRepo, err := repository.NewContentRepository(cfg.BucketName, cfg.S3Region)
	if err != nil {
		log.PrintFatal(ctx, err, nil)
	}

	fetchMetadataService := services.NewFetchMetadataService(metadataRepo, cache, kafkaProducer)
	fetchContentService, err := services.NewFetchContentService(contentRepo, log)
	if err != nil {
		log.PrintFatal(ctx, err, nil)
	}

	return &DownloadCoordinator{
		fetchMetadataService: fetchMetadataService,
		fetchContentService:  fetchContentService,
		cfg:                  cfg,
		logger:               log,
	}, nil
}

func (coord *DownloadCoordinator) DownloadByKey(ctx context.Context, req *pb.DownloadByKeyRequest) (*pb.DownloadByKeyResponse, error) {
	var wg sync.WaitGroup
	errors := make(chan error, 2) // Buffered channel to avoid deadlock

	ress := &pb.DownloadByKeyResponse{}
	wg.Add(2)

	// Fetch metadata
	go func() {
		defer wg.Done()
		metadata, err := coord.fetchMetadataService.FetchMetadataByKey(ctx, req.Key)
		if err != nil {
			errors <- err
			return
		}
		ress.Metadata = metadata
	}()

	// Fetch content
	go func() {
		defer wg.Done()
		url, err := coord.fetchContentService.GetContentUrl(ctx, req.Key)
		if err != nil {
			errors <- err
			return
		}
		ress.DownlaodUrl = url
	}()

	// Wait for goroutines to finish
	wg.Wait()
	close(errors)

	// Check if there were any errors
	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return ress, nil
}

func (coord *DownloadCoordinator) DownloadByUserId(ctx context.Context, req *pb.DownloadByUserIdRequest) (*pb.DownloadByUserIdResponse, error) {
	metadata, err := coord.fetchMetadataService.FetchMetadataByUserId(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	return &pb.DownloadByUserIdResponse{
		Objects: metadata,
	}, nil
}
