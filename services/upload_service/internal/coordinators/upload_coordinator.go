package coordinators

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
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
	var userID *int64
	if req.UserId != "" {
		parsedID, err := strconv.ParseInt(req.UserId, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
		}
		if parsedID <= 0 {
			return nil, status.Error(codes.InvalidArgument, "user ID must be positive integer")
		}
		userID = &parsedID
	}

	metadata := models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
		UserId:         userID,
	}
	if req.UserId == "" {
		metadata.UserId = nil
	} else {
		userId, err := strconv.ParseInt(req.UserId, 0, 16)
		if err != nil {
			return nil, err
		}
		metadata.UserId = &userId
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
