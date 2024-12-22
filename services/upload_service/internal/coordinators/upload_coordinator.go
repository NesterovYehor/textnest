package coordinators

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/api"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
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
	storageRepo, err := repository.NewS3Repository()
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

func (uc *UploadCoordinator) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	metadata := models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errors <- uc.metadataService.ValidateAndSave(ctx, &metadata)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errors <- uc.storageService.SaveContent(ctx, uc.cfg.BucketName, metadata.Key, req.Data)
	}()

	wg.Wait()
	close(errors)

	combinedErr := collectErrors(errors)
	return uc.prepareResponse(metadata.Key, combinedErr)
}

func collectErrors(errors <-chan error) error {
	var errorMessages []string
	for err := range errors {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}
	return nil
}

func (uc *UploadCoordinator) prepareResponse(key string, err error) (*pb.UploadResponse, error) {
	if err != nil {
		uc.log.PrintError(context.Background(), err, map[string]string{"key": key})
		return &pb.UploadResponse{Message: "Failed to upload: " + err.Error()}, nil
	}

	uc.log.PrintInfo(context.Background(), "Upload successful", map[string]string{"key": key})
	return &pb.UploadResponse{Message: "Upload successful"}, nil
}
