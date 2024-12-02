package repository

import (
	"context"

	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
)

type MetadataRepository interface {
	UploadPasteMetadata(ctx context.Context, data *models.MetaData) error
}

type StorageRepository interface {
	UploadPasteContent(
		ctx context.Context,
		bucket, key string,
		data []byte,
	) error
}
