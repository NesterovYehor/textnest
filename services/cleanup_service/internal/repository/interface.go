package repository

import "github.com/NesterovYehor/TextNest/pkg/validator"

type StorageRepository interface {
	DeletePasteContentByKey(key string, bucket string) error
	DeleteExpiredPastes(keys []string, bucket string) error
}

type MetadataRepository interface {
	DeletePasteByKey(key string) error
	DeleteAndReturnExpiredKeys() ([]string, error)
	IsKeyValid(v *validator.Validator, key string)
}
