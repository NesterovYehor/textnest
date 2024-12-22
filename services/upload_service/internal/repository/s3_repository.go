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

type s3Repository struct {
	S3      *s3.Client
	breaker *middleware.CircuitBreakerMiddleware
}

func NewS3Repository() (StorageRepository, error) {
	cbSettings := gobreaker.Settings{
		Name:        "ContentRepo",
		MaxRequests: 1,                // Max requests allowed in half-open state
		Interval:    30 * time.Second, // Time window for tracking errors
		Timeout:     60 * time.Minute, // Time to reset the circuit after tripping
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	log.Printf("Loaded AWS config: %v", cfg.BaseEndpoint)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Ensure virtual-hosted style is used
	})

	return &s3Repository{
		S3:      s3Client,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (repo *s3Repository) UploadPasteContent(ctx context.Context, bucket, key string, data []byte) error {
	_, err := repo.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
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
