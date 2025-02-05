package repository

import (
	"bytes"
	"context"
	"log"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sony/gobreaker"
)

type ContentRepository struct {
	S3      *s3.Client
	bucket  string
	breaker *middleware.CircuitBreakerMiddleware
}

func NewContentRepository(bucket string) (*ContentRepository, error) {
	cbSettings := gobreaker.Settings{
		Name:        "ContentRepo",
		MaxRequests: 3,                // Max requests allowed in half-open state
		Interval:    30 * time.Second, // Time window for tracking errors
		Timeout:     60 * time.Minute, // Time to reset the circuit after tripping
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	log.Printf("Loaded AWS config: %v", cfg.BaseEndpoint)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Ensure virtual-hosted style is used
	})

	return &ContentRepository{
		S3:      s3Client,
		bucket:  bucket,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (repo *ContentRepository) UploadPasteContent(ctx context.Context, key string, data []byte) error {
	_, err := repo.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &repo.bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
		return err
	}
	log.Println("Successfully uploaded object:", key)
	return nil
}
