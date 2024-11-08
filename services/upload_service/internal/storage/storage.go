package storage

import "github.com/NesterovYehor/TextNest/upload_service/models"

type Storage interface {
	GetPaste(key string) (string, error)        // Retrieve post data from blob storage
	UploadPaste(key *models.PasteData) (string, error) // Upload post data to blob storage
	DeletePaste(key string) error
}
