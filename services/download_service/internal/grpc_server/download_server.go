package download_service

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DownloadService struct {
	UnimplementedDownloadServiceServer
	dbRepo      repository.MetadataRepository
	storageRepo repository.StorageRepository
	cfg         *config.Config
}

func NewDownloadService(storageRepo repository.StorageRepository, dbRepo repository.MetadataRepository) *DownloadService {
	return &DownloadService{
		dbRepo:      dbRepo,
		storageRepo: storageRepo,
	}
}

func (srv *DownloadService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error) {
	metadata, err := srv.dbRepo.DownloadPasteMetadata(req.Key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if time.Now().After(metadata.ExpiredDate) {
		producer, err := kafka.NewProducer(*srv.cfg.Kafka, ctx)
		if err != nil {
			return nil, err
		}

		err = producer.ProduceMessages(metadata.Key, srv.cfg.Kafka.Topics["delete-expired-paste"])
		if err != nil {
			return nil, err
		}
	}

	content, err := srv.storageRepo.DownloadPasteContent(srv.cfg.Storage.Bucket, metadata.Key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &DownloadResponse{
		Key:            req.Key,
		ExpirationDate: timestamppb.New(metadata.ExpiredDate),
		CreatedDate:    timestamppb.New(metadata.CreatedAt),
		Content:        content,
	}, nil
}
