package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
)

type metadataRepo struct {
	DB *sql.DB
}

// NewMetadataRepo creates a new instance of MetadataRepository
func NewMetadataRepo(db *sql.DB) MetadataRepository {
	return &metadataRepo{DB: db}
}

// DownloadPasteMetadata retrieves metadata for a paste by key
func (repo *metadataRepo) DownloadPasteMetadata(key string) (*models.Metadata, error) {
	v := validator.New()
	if repo.isKeyValid(key, v); !v.Valid() {
		return nil, errors.New("Key is not valid")
	}

	query := `
        SELECT key, created_at, expired_date FROM metadata WHERE key = $1
    `

	// Set up a context with a timeout for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var paste models.Metadata

	err := repo.DB.QueryRowContext(ctx, query, key).Scan(
		&paste.Key,
		&paste.CreatedAt,
		&paste.ExpiredDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no paste found with the key: %s", key)
		}
		return nil, fmt.Errorf("query failed: %v", err)
	}

	return &paste, nil
}

func (repo *metadataRepo) isKeyValid(key string, v *validator.Validator) {
	v.Check(len([]rune(key)) != 8, "key", "Key must be 8 chars lenth")
}
