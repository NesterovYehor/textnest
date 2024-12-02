package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/sony/gobreaker"

	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
)

type metadataRepository struct {
	DB      *sql.DB
	breaker *middleware.CircuitBreakerMiddleware
}

// NewMetadataRepository creates a new metadata repository with circuit breaker middleware.
func NewMetadataRepository(db *sql.DB) MetadataRepository {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
	}
	return &metadataRepository{
		DB:      db,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}
}

// UploadPasteMetadata inserts metadata into the database with circuit breaker protection.
func (repo *metadataRepository) UploadPasteMetadata(ctx context.Context, data *models.MetaData) error {
	operation := func(ctx context.Context) (any, error) {
		query := `
        INSERT INTO metadata(key, created_at, expiration_date) 
        VALUES ($1, $2, $3)
        `
		// Execute the query
		_, err := repo.DB.ExecContext(ctx, query, data.Key, data.CreatedAt, data.ExpirationDate)
		return nil, err
	}

	// Execute the operation with the circuit breaker
	_, err := repo.breaker.Execute(ctx, operation)
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("request timed out while uploading paste metadata")
	} else if err != nil {
		return err
	}
	return nil
}
