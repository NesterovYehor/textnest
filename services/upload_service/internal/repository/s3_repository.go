package repository

import (
	"context"
	"strings"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sony/gobreaker"
)

type s3Repository struct {
	S3      *s3.Client
	breaker *middleware.CircuitBreakerMiddleware
}

func NewS3Repository(region string) (StorageRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)
	cbSettings := gobreaker.Settings{
		Name:        "ContentRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    30 * time.Second, // Time window for tracking errors
		Timeout:     60 * time.Second, // Time to reset the circuit after tripping
	}

	return &s3Repository{
		S3:      s3Client,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (repo *s3Repository) UploadPasteContent(ctx context.Context, bucket, key string, data []byte) error {
	operation := func(ctx context.Context) (any, error) {
		uploader := manager.NewUploader(repo.S3)

		_, err := uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   strings.NewReader(string(data)),
		})
		return nil, err
	}
	if _, err := repo.breaker.Execute(ctx, operation); err != nil {
		return err
	}
	return nil
}
