package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/models"
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

// GetPaste retrieves a paste from S3 and decodes it into PasteData
func (storage *S3Storage) GetPaste(hash string) (*models.PasteData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Retrieve the object from S3
	res, err := storage.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &storage.Bucket,
		Key:    &hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load data from storage: %w", err)
	}
	defer res.Body.Close()

	var paste models.PasteData

	// Decode the JSON body into PasteData
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&paste); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &paste, nil
}

// UploadPaste uploads paste data to S3 and returns the upload location URL
func (storage *S3Storage) UploadPaste(data *models.PasteData) (string, error) {
	uploader := manager.NewUploader(storage.S3)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Upload the paste data to S3
	res, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &storage.Bucket,
		Key:    &data.Hash,
		Body:   strings.NewReader(data.Content),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload paste: %w", err)
	}

	return res.Location, nil
}

// DeletePaste deletes a paste from S3
func (storage *S3Storage) DeletePaste(hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Delete the object from S3
	_, err := storage.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &storage.Bucket,
		Key:    &hash,
	})
	if err != nil {
		return fmt.Errorf("failed to delete paste: %w", err)
	}
	return nil
}
