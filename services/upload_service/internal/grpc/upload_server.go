package upload_service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/storage"
)

// UploadServer implements the UploadService server.
type UploadServer struct {
	UnimplementedUploadServiceServer // Ensure this is the correct unimplemented server from the generated code
	storage                          storage.Storage
	models                           models.Models
}

// NewUploadServer creates a new instance of UploadServer.
func NewUploadServer(storage storage.Storage, models models.Models) *UploadServer {
	return &UploadServer{
		storage: storage,
		models:  models,
	}
}

// Upload handles the upload request, saving metadata to the database and content to storage.
func (srv *UploadServer) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
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
		if models.IsMetaDataValid(&metadata, v); !v.Valid() {
			for _, err := range v.Errors {
				errCh <- err
			}
		}

		if err := srv.models.MetaData.Insert(&metadata); err != nil {
			errCh <- err.Error()
		} else {
			fmt.Println("Metadata uploaded successfully to DB")
		}
	}()

	// Goroutine to upload the content to storage
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.storage.UploadPaste(metadata.Key, req.Data); err != nil {
			errCh <- err.Error()
		} else {
			fmt.Println("Paste uploaded successfully to storage")
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
		log.Printf("Errors during upload: %v", errorMessages)
		return &UploadResponse{
			Message: "Error(s) occurred: " + strings.Join(errorMessages, ", "),
		}, nil
	}

	return &UploadResponse{
		Message: "Uploaded new paste successfully",
	}, nil
}
