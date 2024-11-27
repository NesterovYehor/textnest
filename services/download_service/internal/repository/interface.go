package repository

import (
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
)

type MetadataRepository interface {
	DownloadPasteMetadata(key string) (*models.Metadata, error)
	IsKeyValid(key string, v *validator.Validator)
}

type StorageRepository interface {
	DownloadPasteContent(bucket, key string) ([]byte, error)
}
