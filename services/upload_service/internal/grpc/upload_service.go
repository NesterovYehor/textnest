package upload_service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/validation"
)

// UploadServer implements the UploadService server.
type UploadService struct {
	UnimplementedUploadServiceServer // Ensure this is the correct unimplemented server from the generated code
	storageRepo                      repository.StorageRepository
	metadataRepo                     repository.MetadataRepository
	cfg                              *config.Config
	log                              *jsonlog.Logger
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
func (srv *UploadService) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	var wg sync.WaitGroup
	metadata := models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
	}
	errCh := make(chan string, 2) // Channel to collect errors from goroutines

	// Goroutine to validate metadata and insert it into the database
	wg.Add(1)
	go func() {
		defer wg.Done()
		v := validator.New()
		if validation.ValidateMetaData(&metadata, v); !v.Valid() {
			for _, err := range v.Errors {
				errCh <- err
			}
		}

		if err := srv.metadataRepo.UploadPasteMetadata(&metadata); err != nil {
			errCh <- err.Error()
		} else {
			srv.log.PrintInfo(ctx, "Metadata uploaded successfully to DB", nil)
		}
	}()

	// Goroutine to upload the content to storage
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.storageRepo.UploadPasteContent(metadata.Key, srv.cfg.BucketName, req.Data); err != nil {
			errCh <- err.Error()
		} else {
			srv.log.PrintInfo(ctx, "Paste uploaded successfully to storage", nil)
		}
	}()

	// Wait for both goroutines to complete and close the error channel
	wg.Wait()
	close(errCh)

	// Collect any errors from the error channel
	var errorMessages []string
	for err := range errCh {
		errorMessages = append(errorMessages, err)
	}

	if len(errorMessages) > 0 {
		// Log and return the collected error messages
		srv.log.PrintError(ctx, fmt.Errorf("Errors during upload: %v", errorMessages), nil)
		return &UploadResponse{
			Message: "Error(s) occurred: " + strings.Join(errorMessages, ", "),
		}, nil
	}

	return &UploadResponse{
		Message: "Uploaded new paste successfully",
	}, nil
}
