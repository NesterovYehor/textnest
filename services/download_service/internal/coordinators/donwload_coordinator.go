package coordinators

import (
	"context"
	"fmt"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DownloadCoordinator struct {
	fetchMetadataService *services.FetchMetadataService
	fetchContentService  *services.FetchContentService
	cfg                  *config.Config
	logger               *jsonlog.Logger
	pb.UnimplementedPasteDownloadServer
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

func (coord *DownloadCoordinator) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	if v := validation.ValidateKey(req.Key); !v.Valid() {
		return nil, fmt.Errorf("invalid key")
	}

	// Fetch metadata
	metadata, err := coord.fetchMetadataService.GetMetadata(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Fetch content
	content, err := coord.fetchContentService.GetContent(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Build and return the response
	return coord.createDownloadResponse(req.Key, models.FromProto(metadata), content), nil
}

func (coord *DownloadCoordinator) createDownloadResponse(key string, metadata *models.Metadata, content []byte) *pb.DownloadResponse {
	return &pb.DownloadResponse{
		Key:            key,
		ExpirationDate: timestamppb.New(metadata.ExpiredDate),
		CreatedDate:    timestamppb.New(metadata.CreatedAt),
		Content:        content,
	}
}
