package repository

import (
	"context"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sony/gobreaker"
)

type ContentRepo struct {
	S3     *s3.PresignClient
	beaker *middleware.CircuitBreakerMiddleware
	bucket string
}

func NewContentRepository(bucket, region string) (*ContentRepo, error) {
	// Load AWS configuration with specified region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    10 * time.Second, // Time window for tracking errors
		Timeout:     30 * time.Second, // Time to reset the circuit after tripping
	}

	presignClient := s3.NewPresignClient(s3.NewFromConfig(cfg))

	return &ContentRepo{
		S3:     presignClient,
		beaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
		bucket: bucket,
	}, nil
}

func (repo *ContentRepo) GenerateDownloadURL(key string, ctx context.Context) (string, error) {
	req, err := repo.S3.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &repo.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(time.Minute*10))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

