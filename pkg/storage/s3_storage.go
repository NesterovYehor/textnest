package storage

import (
	"context"
	"fmt"
	"io"
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

// GetPaste retrieves a paste from S3 and decodes it into PasteData.
func (storage *S3Storage) GetPaste(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Retrieve the object from S3 using the bucket name
	res, err := storage.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(storage.Bucket), // Use the correct bucket name
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load data from storage: %w", err)
	}
	defer res.Body.Close()

	// Read the object data
	pasteContent, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	// Convert the byte slice to a string and return it
	return pasteContent, nil
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

// DeletePaste deletes a paste from S3
func (storage *S3Storage) DeletePaste(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Delete the object from S3
	_, err := storage.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(storage.Bucket), // Use the correct bucket name
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete paste: %w", err)
	}
	return nil
}
