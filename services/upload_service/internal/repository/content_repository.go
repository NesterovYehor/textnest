package repository

import (
	"context"
	"log"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sony/gobreaker"
)

type ContentRepository struct {
	client  *s3.PresignClient
	bucket  string
	breaker *middleware.CircuitBreakerMiddleware
}

func NewContentRepository(bucket, region string) (*ContentRepository, error) {
	cbSettings := gobreaker.Settings{
		Name:        "ContentRepo",
		MaxRequests: 3,                // Max requests allowed in half-open state
		Interval:    30 * time.Second, // Time window for tracking errors
		Timeout:     60 * time.Minute, // Time to reset the circuit after tripping
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewPresignClient(s3.NewFromConfig(cfg))

	return &ContentRepository{
		client:  client,
		bucket:  bucket,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (repo *ContentRepository) GenerateUploadURL(ctx context.Context, key string) (string, error) {
	req, err := repo.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &repo.bucket,
		Key:    &key,
	})
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
		return "", err
	}
	return req.URL, nil
}
