package download_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DownloadServer struct {
	UnimplementedDownloadServiceServer
	storage storage.Storage
	model   models.Models
}

func NewDownloadServer(storage storage.Storage, model models.Models) *DownloadServer {
	return &DownloadServer{
		storage: storage,
		model:   model,
	}
}

func (srv *DownloadServer) Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error) {
	v := validator.New()

	if srv.model.Metadata.IsKeyValid(req.Key, v); v.Valid() {
		return nil, errors.New("Key is invalid")
	}

	metadata, err := srv.model.Metadata.Get(req.Key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	content, err := srv.storage.DownloadPaste(req.Key)
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
