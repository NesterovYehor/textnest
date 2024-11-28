package repository

import (
	"context"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sony/gobreaker"
)

type s3Repository struct {
	S3     *s3.Client
	beaker *middleware.CircuitBreakerMiddleware
}

func NewStorageRepository(bucket, region string) (StorageRepository, error) {
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

	// Initialize the S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &s3Repository{
		S3:     s3Client,
		beaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (repo *s3Repository) DownloadPasteContent(bucket, key string) ([]byte, error) {
	operation := func(ctx context.Context) (any, error) {
		// Create the downloader with the S3 client
		downloader := manager.NewDownloader(repo.S3)

		// Set a timeout for the context
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Use manager.WriteAtBuffer to store the downloaded content
		buf := manager.NewWriteAtBuffer([]byte{})

		// Create the GetObjectInput with the S3 key and bucket name
		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}

		// Perform the download operation and capture the result into the buffer
		_, err := downloader.Download(ctx, buf, input)
		if err != nil {
			return nil, fmt.Errorf("failed to download object from S3: %v", err)
		}

		// Return the content as a byte slice
		return buf.Bytes(), nil
	}

	// Use the Circuit Breaker middleware to execute the operation
	result, err := repo.beaker.Execute(context.Background(), operation)
	if err != nil {
		return nil, err
	}

	// Ensure the result is cast correctly to a byte slice
	content, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from Circuit Breaker execution")
	}

	return content, nil
}
