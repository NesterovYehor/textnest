package repository

import (

	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
)

type MetadataRepository interface {
	DownloadPasteMetadata(key string) (*models.Metadata, error)
}

type StorageRepository interface {
	DownloadPasteContent(bucket, key string) ([]byte, error)
}
