package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/sony/gobreaker"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (repo *MetadataRepo) DownloadPasteMetadata(ctx context.Context, key string) (*pb.Metadata, error) {
	operation := func(ctx context.Context) (any, error) {
		query := `SELECT key, title, created_at, expiration_date FROM metadata WHERE key = $1`
		var paste pb.Metadata
		var createdAt time.Time
		var expiredDate time.Time

		err := repo.DB.QueryRowContext(ctx, query, key).Scan(
			&paste.Key,
			&paste.Title,
			&createdAt,
			&expiredDate,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("no paste found with key %s: %w", key, err)
			}
			return nil, fmt.Errorf("query failed: %w", err)
		}
		paste.CreatedAt = timestamppb.New(createdAt)
		paste.ExpiredDate = timestamppb.New(expiredDate)
		return &paste, nil
	}

	result, err := repo.breaker.Execute(ctx, operation)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker error: %w", err)
	}

	paste, ok := result.(*pb.Metadata)
	if !ok {
		return nil, errors.New("unexpected result type")
	}
	return paste, nil
}

func (repo *MetadataRepo) DownloadMetadataByUserId(ctx context.Context, userId string, limit, offset int) ([]*pb.Metadata, error) {
	operation := func(ctx context.Context) (any, error) {
		query := `SELECT key, title, created_at, expiration_date FROM metadata WHERE user_id = $1 LIMIT $2 OFFSET $3`
		rows, err := repo.DB.QueryContext(ctx, query, userId, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		var metadata []*pb.Metadata
		for rows.Next() {
			var expiredDate time.Time
			var createdAt time.Time
			var m pb.Metadata
			if err := rows.Scan(&m.Key, &m.Title, &createdAt, &expiredDate); err != nil { // FIX: Added `&` before m.Title
				return nil, fmt.Errorf("scan failed: %w", err)
			}
			m.CreatedAt = timestamppb.New(createdAt)
			m.ExpiredDate = timestamppb.New(expiredDate)
			metadata = append(metadata, &m)
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

	metadata, ok := result.([]*pb.Metadata)
	if !ok {
		return nil, errors.New("unexpected result type")
	}
	return metadata, nil
}
