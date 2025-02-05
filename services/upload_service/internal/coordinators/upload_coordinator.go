package coordinators

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UploadCoordinator struct {
	metadataService *services.MetadataManagementService
	storageService  *services.ContentManagementService
	cfg             *config.Config
	log             *jsonlog.Logger
	pb.UnimplementedPasteUploadServer
}

func NewUploadCoordinator(cfg *config.Config, log *jsonlog.Logger, db *sql.DB) (*UploadCoordinator, error) {
	metadataRepo := repository.NewMetadataRepository(db)
	storageRepo, err := repository.NewContentRepository(cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 repository: %w", err)
	}
	return &UploadCoordinator{
		metadataService: services.NewMetadataManagementService(metadataRepo, log),
		storageService:  services.NewStorageService(storageRepo, log),
		cfg:             cfg,
		log:             log,
	}, nil
}

func (uc *UploadCoordinator) UploadPaste(ctx context.Context, req *pb.UploadPasteRequest) (*pb.UploadPasteResponse, error) {
	metadata := models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
		UserId:         req.UserId,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate and save metadata first
	if err := uc.metadataService.ValidateAndSave(ctx, &metadata); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save metadata: %v", err)
	}

	// Save content with cleanup on failure
	if err := uc.storageService.SaveContent(ctx, metadata.Key, req.Data); err != nil {
		// Attempt to mark as failed, log but ignore secondary errors
		if expireErr := uc.metadataService.ExpireMetadata(ctx, metadata.Key); expireErr != nil {
			uc.log.PrintError(ctx, fmt.Errorf("failed to mark paste as failed"),
				map[string]string{
					"key":   metadata.Key,
					"error": expireErr.Error(),
				})
		}

		return nil, status.Errorf(codes.Internal, "content save failed: %v", err)
	}

	uc.log.PrintInfo(ctx, "upload successful", map[string]string{"key": metadata.Key})
	return &pb.UploadPasteResponse{Message: "Upload successful"}, nil
}

func (uc *UploadCoordinator) UploadUpdates(ctx context.Context, req *pb.UploadUpdatesRequest) (*pb.UploadUpdatesResponse, error) {
	userId, err := uc.metadataService.GetPasteOwner(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if userId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "Upload content failed %v", err)
	}

	var wg sync.WaitGroup
	errorCh := make(chan error, 2)

	if req.ExpirationDate != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errorCh <- uc.metadataService.UpdateMetadata(ctx, req.Key, req.ExpirationDate.AsTime())
		}()
	}
	if req.Content != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errorCh <- uc.storageService.UpdateContent(ctx, req.Key, req.Content)
		}()
	}

	wg.Wait()
	close(errorCh)

	combineErr := collectErrors(errorCh)
	if combineErr != nil {
		return nil, combineErr
	}
	return &pb.UploadUpdatesResponse{Message: "Paste updated successfully"}, nil
}

func (uc *UploadCoordinator) ExpirePaste(ctx context.Context, req *pb.ExpirePasteRequest) (*pb.ExpirePasteResponse, error) {
	userId, err := uc.metadataService.GetPasteOwner(ctx, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Paste is not found")
	}
	if req.UserId != userId {
		return nil, status.Errorf(codes.PermissionDenied, "only author can expire paste")
	}
	if err := uc.metadataService.UpdateMetadata(ctx, req.Key, time.Now()); err != nil {
		return nil, err
	}
	return &pb.ExpirePasteResponse{
		Message: "Paste expired successfully",
	}, nil
}

func (uc *UploadCoordinator) ExpireAllPastesByUserID(ctx context.Context, req *pb.ExpireAllPastesByUserIDRequest) (*pb.ExpireAllPastesByUserIDResponse, error) {
	if err := uc.metadataService.ExpireAllPastes(ctx, req.UserId); err != nil {
		return nil, err
	}
	return &pb.ExpireAllPastesByUserIDResponse{Message: fmt.Sprintf("All Pastes of user %v expired successfully", req.UserId)}, nil
}

func collectErrors(errorCh <-chan error) error {
	var errorMessages []string
	for err := range errorCh {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}
	return nil
}
