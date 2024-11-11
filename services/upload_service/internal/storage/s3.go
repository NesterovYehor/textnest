package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Bucket string
	S3     *s3.Client
}

func NewS3Storage(bucket, region string) (*S3Storage, error) {
	// Load AWS configuration with specified region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &S3Storage{
		Bucket: bucket,
		S3:     s3Client,
	}, nil
}

// UploadPaste uploads paste data to S3 and returns the upload location URL
func (storage *S3Storage) UploadPaste(key string, data []byte) error {
	uploader := manager.NewUploader(storage.S3)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println(storage.Bucket)

	// Upload the paste data to S3
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(storage.Bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(string(data)),
		ACL:    "public-read",
	})
	if err != nil {
		return err
	}

	return nil
}
