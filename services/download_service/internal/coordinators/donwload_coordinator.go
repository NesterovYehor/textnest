package coordinators

import (
	"context"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
)

type DownloadCoordinator struct {
	fetchMetadataService *services.FetchMetadataService
	fetchContentService  *services.FetchContentService
	cfg                  *config.Config
	logger               *jsonlog.Logger
	pb.UnsafePasteDownloadServer
}

func NewDownloadCoordinator(ctx context.Context, cfg *config.Config, log *jsonlog.Logger) (*DownloadCoordinator, error) {
	factory := services.NewServiceFactory(cfg, log)
	fetchMetadataService, err := factory.CreateFetchMetadataService(ctx)
	if err != nil {
		return nil, err
	}
	fetchContentService, err := factory.CreateFetchContentService(ctx)
	if err != nil {
		return nil, err
	}

	return &DownloadCoordinator{
		fetchMetadataService: fetchMetadataService,
		fetchContentService:  fetchContentService,
		cfg:                  cfg,
		logger:               log,
	}, nil
}

func (coord *DownloadCoordinator) DownloadByKey(ctx context.Context, req *pb.DownloadByKeyRequest) (*pb.DownloadByKeyResponse, error) {
	// Fetch metadata
	metadata, err := coord.fetchMetadataService.FetchMetadataByKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Fetch content
	content, err := coord.fetchContentService.GetContent(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Build and return the response
	return &pb.DownloadByKeyResponse{
		Key:            req.Key,
		ExpirationDate: metadata.ExpiredDate,
		CreatedDate:    metadata.CreatedAt,
		Content:        content,
	}, nil
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
