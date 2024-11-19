package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Repository struct {
	S3 *s3.Client
}

func NewStorageRepository(bucket, region string) (StorageRepository, error) {
	// Load AWS configuration with specified region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &s3Repository{
		S3: s3Client,
	}, nil
}

func (repo *s3Repository) DownloadPasteContent(bucket, key string) ([]byte, error) {
	// Create the downloader with the S3 client
	downloader := manager.NewDownloader(repo.S3)

	// Set a timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
