package app

import (
	"context"
	"database/sql"
	"sync"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	grpc_clients "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client"
)

type AppContext struct {
	DB           *sql.DB
	KeyGenClient *grpc_clients.KeyGeneratorClient
	UploadClient *grpc_clients.UploadClient
	Logger       *jsonlog.Logger
}

var (
	instance *AppContext
	once     sync.Once
)

// GetAppContext initializes and returns the singleton instance of AppContext.
func GetAppContext(cfg *config.Config, ctx context.Context, logger *jsonlog.Logger) (*AppContext, error) {
	var err error
	once.Do(func() {
		// Create a gRPC client for Key Generation Service
		keyGenClient, keyGenErr := grpc_clients.NewKeyGeneratorClient(cfg.KeyService.Port)
		if keyGenErr != nil {
			err = keyGenErr
			return
		}

		// Create a gRPC client for Upload Service
		uploadPasteClient, uploadErr := grpc_clients.NewUploadClient(cfg.UploadService.Port)
		if uploadErr != nil {
			err = uploadErr
			return
		}
		// Set the singleton instance
		instance = &AppContext{
			KeyGenClient: keyGenClient,
			UploadClient: uploadPasteClient,
			Logger:       logger,
		}
	})

	return instance, err
}

// Close releases resources held by the AppContext
func (a *AppContext) Close() {
	if a.DB != nil {
		_ = a.DB.Close()
	}

	// Note: gRPC connections can also be closed if explicitly stored.
}

