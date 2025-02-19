package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
	"github.com/sony/gobreaker"
)

type MetadataRepository struct {
	DB      *sql.DB
	breaker *middleware.CircuitBreakerMiddleware
}

// NewMetadataRepository creates a new metadata repository with circuit breaker middleware.
func NewMetadataRepository(db *sql.DB) *MetadataRepository {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	}
	return &MetadataRepository{
		DB:      db,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}
}

// UploadPasteMetadata inserts metadata into the database with circuit breaker protection.
func (repo *MetadataRepository) InsertPasteMetadata(ctx context.Context, data *pb.UploadPasteRequest) error {
	operation := func(ctx context.Context) (any, error) {
		query := `
        INSERT INTO metadata(key, title, user_id, expiration_date) 
        VALUES ($1, NULLIF($2, ''), $3, $4)
        `

		args := []any{
			data.Key,
            data.Title,
			data.UserId,
			data.ExpirationDate.AsTime(),
		}
		// Execute the query
		_, err := repo.DB.ExecContext(ctx, query, args...)
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

func (repo *MetadataRepository) UpdatePasteMetadata(ctx context.Context, expirationDate time.Time, key string) error {
	operation := func(ctx context.Context) (any, error) {
		// Try to update the metadata
		query := `
        UPDATE metadata 
        SET expiration_date = $1
        WHERE key = $2
        `
		args := []any{
			expirationDate,
			key,
		}
		_, err := repo.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Execute the operation with the circuit breaker
	_, err := repo.breaker.Execute(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MetadataRepository) ExpireAllPastesByUserId(ctx context.Context, userId string) error {
	operation := func(ctx context.Context) (any, error) {
		query := `UPDATE metadata SET expiration_date = NOW() WHERE user_id = $1`
		_, err := repo.DB.ExecContext(ctx, query, userId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := repo.breaker.Execute(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MetadataRepository) GetPasteOwner(ctx context.Context, key string) (string, error) {
	query := `SELECT user_id FROM metadata WHERE key = $1`

	var userId string
	err := repo.DB.QueryRowContext(ctx, query, key).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
