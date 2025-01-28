package app

import (
	"context"
	"fmt"
	"sync"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	grpc_clients "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client"
)

type AppContext struct {
	KeyGenClient   *grpc_clients.KeyGeneratorClient
	UploadClient   *grpc_clients.UploadClient
	AuthClient     *grpc_clients.AuthClient
	DownloadClient *grpc_clients.DownloadClient
	closers        []func() error
	Logger         *jsonlog.Logger
}

var (
	instance *AppContext
	once     sync.Once
)

// GetAppContext initializes and returns the singleton instance of AppContext.
func GetAppContext(cfg *config.Config, ctx context.Context, logger *jsonlog.Logger) (*AppContext, error) {
	var err error
	once.Do(func() {
		instance = &AppContext{
			Logger: logger,
		}

		// Create a gRPC client for Key Generation Service
		keyGenClient, keyGenErr := grpc_clients.NewKeyGeneratorClient(cfg.KeyService.Port)
		if keyGenErr != nil {
			err = keyGenErr
			return
		}
		instance.closers = append(instance.closers, keyGenClient.Close)

		// Create a gRPC client for Upload Service
		uploadPasteClient, uploadErr := grpc_clients.NewUploadClient(cfg.UploadService.Port)
		if uploadErr != nil {
			err = uploadErr
			return
		}
		instance.closers = append(instance.closers, keyGenClient.Close)
		downloadPasteClient, downloadErr := grpc_clients.NewDownloadClient(cfg.DownloadService.Port)
		if downloadErr != nil {
			err = uploadErr
			return
		}

		instance.closers = append(instance.closers, downloadPasteClient.Close)
		authClient, authErr := grpc_clients.NewAuthClient(cfg.AuthService.Port)
		if authErr != nil {
			err = authErr
			return
		}
		instance.closers = append(instance.closers, authClient.Close)
		// Set the singleton instance
		instance = &AppContext{
			KeyGenClient:   keyGenClient,
			UploadClient:   uploadPasteClient,
			DownloadClient: downloadPasteClient,
			AuthClient:     authClient,
			Logger:         logger,
		}
	})

	return instance, err
}

func (app *AppContext) Close() error {
	var errs []error
	for _, close := range app.closers {
		if err := close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing resources: %v", errs)
	}
	return nil
}
