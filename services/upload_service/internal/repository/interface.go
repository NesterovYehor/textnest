package repository

import "github.com/NesterovYehor/TextNest/services/upload_service/internal/models"

type MetadataRepository interface {
	UploadPasteMetadata(data *models.MetaData) error
}

type StorageRepository interface {
	UploadPasteContent(bucket, key string, data []byte) error
}
