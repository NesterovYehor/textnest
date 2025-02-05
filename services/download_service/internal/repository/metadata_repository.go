package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/models"
	"github.com/sony/gobreaker"
)

type MetadataRepo struct {
	DB      *sql.DB
	breaker *middleware.CircuitBreakerMiddleware
}

// NewMetadataRepo creates a new instance of MetadataRepository
func NewMetadataRepo(db *sql.DB) *MetadataRepo {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    10 * time.Second, // Time window for tracking errors
		Timeout:     30 * time.Second, // Time to reset the circuit after tripping
	}
	return &MetadataRepo{
		DB:      db,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}
}

func (repo *MetadataRepo) DownloadPasteMetadata(ctx context.Context, key string) (*models.Metadata, error) {
	operation := func(ctx context.Context) (any, error) {
		query := `SELECT key, created_at, expiration_date FROM metadata WHERE key = $1`
		var paste models.Metadata

		err := repo.DB.QueryRowContext(ctx, query, key).Scan(
			&paste.Key,
			&paste.CreatedAt,
			&paste.ExpirationDate, // Corrected field name
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("no paste found with key %s: %w", key, err)
			}
			return nil, fmt.Errorf("query failed: %w", err)
		}
		return &paste, nil
	}

	result, err := repo.breaker.Execute(ctx, operation)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker error: %w", err)
	}

	paste, ok := result.(*models.Metadata)
	if !ok {
		return nil, errors.New("unexpected result type")
	}
	return paste, nil
}

func (repo *MetadataRepo) DownloadMetadataByUserId(ctx context.Context, userId string, limit, offcet int) ([]models.Metadata, error) {
	operation := func(ctx context.Context) (any, error) {
        query := `SELECT key, created_at, expiration_date FROM metadata WHERE user_id = $1 LIMIT $2 OFFSET $3`
		rows, err := repo.DB.QueryContext(ctx, query, userId, limit, offcet)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		var metadata []models.Metadata
		for rows.Next() {
			var m models.Metadata
			if err := rows.Scan(&m.Key, &m.CreatedAt, &m.ExpirationDate); err != nil {
				return nil, fmt.Errorf("scan failed: %w", err)
			}
			metadata = append(metadata, m)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows error: %w", err)
		}
		return metadata, nil
	}

	result, err := repo.breaker.Execute(ctx, operation)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker error: %w", err)
	}

	metadata, ok := result.([]models.Metadata)
	if !ok {
		return nil, errors.New("unexpected result type")
	}
	return metadata, nil
}
