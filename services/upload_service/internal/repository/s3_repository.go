package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Repository struct {
	S3 *s3.Client
}

func NewS3Repository(region string) (StorageRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)

	return &s3Repository{
		S3: s3Client,
	}, nil
}

func (repo *s3Repository) UploadPasteContent(bucket, key string, data []byte) error {
	uploader := manager.NewUploader(repo.S3)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(string(data)),
	})
	if err != nil {
		return err
	}

	return nil
}
