package coordinators

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UploadCoordinator struct {
	metadataService *services.MetadataManagementService
	storageService  *services.ContentManagementService
	mu              sync.Mutex
	cfg             *config.Config
	log             *jsonlog.Logger
	pb.UnimplementedPasteUploadServer
}

func NewUploadCoordinator(cfg *config.Config, log *jsonlog.Logger, db *sql.DB) (*UploadCoordinator, error) {
	metadataRepo := repository.NewMetadataRepository(db)
	storageRepo, err := repository.NewContentRepository(cfg.BucketName, cfg.S3Region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 repository: %w", err)
	}
	return &UploadCoordinator{
		metadataService: services.NewMetadataManagementService(metadataRepo, log),
		storageService:  services.NewStorageService(storageRepo, log),
		mu:              sync.Mutex{},
		cfg:             cfg,
		log:             log,
	}, nil
}

func (uc *UploadCoordinator) UploadPaste(ctx context.Context, req *pb.UploadPasteRequest) (*pb.UploadPasteResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		resp    pb.UploadPasteResponse
		errChan = make(chan error, 2)
		urlChan = make(chan string, 1)
		wg      sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := uc.metadataService.ValidateAndSave(ctx, req); err != nil {
			errChan <- fmt.Errorf("metadata save: %w", err)
			cancel()
		} else {
			resp.ExpirationDate = timestamppb.New(req.ExpirationDate.AsTime())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		url, err := uc.storageService.GenerateUploadURL(ctx, req.Key)
		if err != nil {
			errChan <- fmt.Errorf("generate URL: %w", err)
			cancel()
			return
		}
		urlChan <- url
	}()

	wg.Wait()
	close(errChan)
	close(urlChan)

	if err := collectErrors(errChan); err != nil {
		_ = uc.metadataService.ExpireMetadata(context.Background(), req.Key)
		return nil, status.Errorf(codes.Internal, "upload failed: %v", err)
	}

	resp.UploadUrl = <-urlChan
	return &resp, nil
}

func (uc *UploadCoordinator) UploadUpdates(ctx context.Context, req *pb.UploadUpdatesRequest) (*pb.UploadUpdatesResponse, error) {
	userId, err := uc.metadataService.GetPasteOwner(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if userId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "Upload content failed %v", err)
	}
	url, err := uc.storageService.GenerateUploadURL(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	return &pb.UploadUpdatesResponse{UploadUrl: url}, nil
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
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}
