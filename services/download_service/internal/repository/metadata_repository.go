package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	middleware "github.com/NesterovYehor/TextNest/services/download_service/internal/middlewares"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/sony/gobreaker"
)

type metadataRepo struct {
	DB      *sql.DB
	breaker *middleware.CircuitBreakerMiddleware
}

// NewMetadataRepo creates a new instance of MetadataRepository
func NewMetadataRepo(db *sql.DB) MetadataRepository {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    10 * time.Second, // Time window for tracking errors
		Timeout:     30 * time.Second, // Time to reset the circuit after tripping
	}
	return &metadataRepo{
		DB:      db,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}
}

// DownloadPasteMetadata retrieves metadata for a paste by key, with Circuit Breaker protection
func (repo *metadataRepo) DownloadPasteMetadata(key string) (*models.Metadata, error) {
	// Define the operation to be executed with Circuit Breaker
	operation := func(ctx context.Context) (any, error) {
		query := `
            SELECT key, created_at, expired_date FROM metadata WHERE key = $1
        `
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

	// Use the Circuit Breaker middleware to execute the operation
	result, err := repo.breaker.Execute(context.Background(), operation)
	if err != nil {
		return nil, err // Circuit Breaker error or underlying error
	}

	// Cast result to the expected type
	paste, ok := result.(*models.Metadata)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	return paste, nil
}

func (repo *metadataRepo) IsKeyValid(key string, v *validator.Validator) {
	v.Check(len([]rune(key)) != 8, "key", "Key must be 8 chars lenth")
}
