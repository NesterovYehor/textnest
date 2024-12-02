package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/validation"
)

// UploadServer implements the UploadService server.
type UploadService struct {
	pb.UnimplementedUploadServiceServer // Ensure this is the correct unimplemented server from the generated code
	storageRepo                         repository.StorageRepository
	metadataRepo                        repository.MetadataRepository
	cfg                                 *config.Config
	log                                 *jsonlog.Logger
}

// NewUploadServer creates a new instance of UploadServer.
func NewUploadService(storageRepo repository.StorageRepository, metadataRepo repository.MetadataRepository, log *jsonlog.Logger, cfg *config.Config) *UploadService {
	return &UploadService{
		storageRepo:  storageRepo,
		metadataRepo: metadataRepo,
		cfg:          cfg,
		log:          log,
	}
}

// Upload handles the upload request, saving metadata to the database and content to storage.
func (srv *UploadService) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	metadata := srv.prepareMetadata(req)

	// Use channels to capture errors from goroutines
	metadataErrCh := make(chan error, 1)
	storageErrCh := make(chan error, 1)

	// Start goroutines for metadata validation/storage and content upload
	go func() {
		metadataErrCh <- srv.validateAndSaveMetadata(ctx, &metadata)
	}()
	go func() {
		storageErrCh <- srv.saveContentToStorage(ctx, metadata.Key, req.Data)
	}()

	// Collect errors from channels
	validationErr := <-metadataErrCh
	storageErr := <-storageErrCh

	// Combine errors and prepare response
	return srv.prepareResponse(ctx, metadata.Key, validationErr, storageErr)
}

// prepareMetadata creates a MetaData object from the UploadRequest.
func (srv *UploadService) prepareMetadata(req *pb.UploadRequest) models.MetaData {
	return models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
	}
}

// validateAndSaveMetadata validates metadata and saves it to the database.
func (srv *UploadService) validateAndSaveMetadata(ctx context.Context, metadata *models.MetaData) error {
	v := validator.New()
	if validation.ValidateMetaData(metadata, v); !v.Valid() {
		validationErr := fmt.Errorf("metadata validation errors: %v", v.Errors)
		srv.log.PrintError(ctx, validationErr, map[string]string{"key": metadata.Key})
		return validationErr
	}

	if err := srv.metadataRepo.UploadPasteMetadata(ctx, metadata); err != nil {
		dbErr := fmt.Errorf("failed to upload metadata: %w", err)
		srv.log.PrintError(ctx, dbErr, map[string]string{"key": metadata.Key})
		return dbErr
	}

	srv.log.PrintInfo(ctx, "Metadata uploaded successfully to DB", map[string]string{"key": metadata.Key})
	return nil
}

// saveContentToStorage uploads the paste content to storage.
func (srv *UploadService) saveContentToStorage(ctx context.Context, key string, data []byte) error {
	if err := srv.storageRepo.UploadPasteContent(ctx, srv.cfg.BucketName, key, data); err != nil {
		storageErr := fmt.Errorf("failed to upload content: %w", err)
		srv.log.PrintError(ctx, storageErr, map[string]string{"key": key})
		return storageErr
	}

	srv.log.PrintInfo(ctx, "Paste uploaded successfully to storage", map[string]string{"key": key})
	return nil
}

// prepareResponse combines validation and storage errors into a response.
func (srv *UploadService) prepareResponse(ctx context.Context, key string, validationErr, storageErr error) (*pb.UploadResponse, error) {
	if validationErr != nil || storageErr != nil {
		errorMessages := []string{}
		if validationErr != nil {
			errorMessages = append(errorMessages, validationErr.Error())
		}
		if storageErr != nil {
			errorMessages = append(errorMessages, storageErr.Error())
		}

		srv.log.PrintError(ctx, fmt.Errorf("errors occurred during upload"), map[string]string{
			"key":    key,
			"errors": strings.Join(errorMessages, "; "),
		})

		return &pb.UploadResponse{
			Message: "Error(s) occurred: " + strings.Join(errorMessages, "; "),
		}, nil
	}

	srv.log.PrintInfo(ctx, "Uploaded new paste successfully", map[string]string{"key": key})
	return &pb.UploadResponse{
		Message: "Uploaded new paste successfully",
	}, nil
}
