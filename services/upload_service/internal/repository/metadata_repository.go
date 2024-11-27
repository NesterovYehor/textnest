package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
)

type metadataRepository struct {
	DB *sql.DB
}

func NewMetadataRepository(db *sql.DB) MetadataRepository {
	return &metadataRepository{DB: db}
}

func (repo *metadataRepository) UploadPasteMetadata(data *models.MetaData) error {
	query := `
        INSERT INTO metadata(key, created_at, expiration_date) 
        VALUES ($1, $2, $3)
    `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, data.Key, data.CreatedAt, data.ExpirationDate)
	if err != nil {
		return err
	}

	return nil
}
