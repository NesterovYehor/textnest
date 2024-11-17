package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	Bucket            string
	S3                *s3.Client
	s3PresignedClient *s3.PresignClient
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

func (storage *S3Storage) DeletePaste(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := storage.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(storage.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to delete object %v:%v. Error: %v\n", storage.Bucket, key, err)
		return err
	}

	return nil
}

func (storage *S3Storage) DeleteExpiredPastes(keys []string) error {
	if len(keys) == 0 {
		log.Println("No keys provided for deletion.")
		return nil
	}

	var objects []types.ObjectIdentifier
	for _, key := range keys {
		objects = append(objects, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	output, err := storage.S3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(storage.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		log.Printf("Failed to delete objects from bucket %s. Error: %v", storage.Bucket, err)
		return err
	}

	// Log deleted objects for transparency
	if output.Deleted != nil {
		for _, deleted := range output.Deleted {
			log.Printf("Deleted object: %s", *deleted.Key)
		}
	}

	return nil
}
