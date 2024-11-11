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

type UploadServer struct {
	UnimplementedUploadServiceServer
	storage storage.Storage
	models  models.Models
}

func NewUploadServer(storage storage.Storage, models models.Models) *UploadServer {
	return &UploadServer{
		storage: storage,
		models:  models,
	}
}

func (srv *UploadServer) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	var wg sync.WaitGroup
	metadata := models.MetaData{
		Key:            req.Key,
		ExpirationDate: req.ExpirationDate.AsTime(),
		CreatedAt:      time.Now(),
	}
	errCh := make(chan string, 2) // Buffered channel to collect errors

	// Goroutine to validate metadata and insert into DB
	wg.Add(1)
	go func() {
		defer wg.Done()
		v := validator.New()
		if models.IsMetaDataValid(&metadata, v); !v.Valid() {
			for _, err := range v.Errors {
				errCh <- err
			}
		}

		err := srv.models.Paste.Insert(&metadata)
		if err != nil {
			errCh <- err.Error()
		}
		fmt.Println("Metadata uploaded sucesfully to db")
	}()

	// Goroutine to handle the upload to storage
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.storage.UploadPaste(metadata.Key, req.Data)
		if err != nil {
			errCh <- err.Error()
		}
		fmt.Println("Paste uploaded sucesfully to storage")
	}()
	// Wait for both goroutines to finish and check for errors
	wg.Wait()
	close(errCh) // Close the channel after both goroutines are done

	// Collect and handle any errors
	var errorMessages []string
	for err := range errCh {
		errorMessages = append(errorMessages, err)
	}

	if len(errorMessages) > 0 {
		// Return all collected error messages
		log.Printf("Errors during upload: %v", errorMessages)
		return &UploadResponse{
			Message: "Error(s) occurred: " + strings.Join(errorMessages, ", "),
		}, nil
	}

	return &UploadResponse{
		Message: "Uploaded new paste successfully",
	}, nil
}
