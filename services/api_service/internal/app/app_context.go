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

// GetAppContext initializes the singleton AppContext instance.
func GetAppContext(cfg *config.Config, ctx context.Context, logger *jsonlog.Logger) (*AppContext, error) {
	var err error
	once.Do(func() {
		instance = &AppContext{Logger: logger}

		// Create gRPC clients
		keyGenClient, keyGenErr := grpc_clients.NewKeyGeneratorClient(cfg.KeyService.Port)
		if keyGenErr != nil {
			err = keyGenErr
			return
		}
		instance.KeyGenClient = keyGenClient
		instance.closers = append(instance.closers, keyGenClient.Close)

		uploadPasteClient, uploadErr := grpc_clients.NewUploadClient(cfg.UploadService.Port)
		if uploadErr != nil {
			err = uploadErr
			return
		}
		instance.UploadClient = uploadPasteClient
		instance.closers = append(instance.closers, uploadPasteClient.Close)

		downloadPasteClient, downloadErr := grpc_clients.NewDownloadClient(cfg.DownloadService.Port)
		if downloadErr != nil {
			err = downloadErr // FIX: Assign correct error variable
			return
		}
		instance.DownloadClient = downloadPasteClient
		instance.closers = append(instance.closers, downloadPasteClient.Close)

		authClient, authErr := grpc_clients.NewAuthClient(cfg.AuthService.Port)
		if authErr != nil {
			err = authErr
			return
		}
		instance.AuthClient = authClient
		instance.closers = append(instance.closers, authClient.Close)
	})

	return instance, err
}

// GetInstance returns the existing AppContext instance without reinitialization.
func GetInstance() (*AppContext, error) {
	if instance == nil {
		return nil, fmt.Errorf("app context is not initialized, call GetAppContext first")
	}
	return instance, nil
}

// Close releases resources.
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

