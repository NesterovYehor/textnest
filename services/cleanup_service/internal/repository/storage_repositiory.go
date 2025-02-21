package repository

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

type StorageRepo struct {
	bucketName string
	client         *s3.Client
}

func NewStorageRepo(region, bucket string) (*StorageRepo, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &StorageRepo{
		client:         s3Client,
		bucketName: bucket,
	}, nil
}

func (storage *StorageRepo) DeletePasteContentByKey(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := storage.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(storage.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to delete object %v:%v. Error: %v\n", storage.bucketName, key, err)
		return err
	}

	return nil
}

func (storage *StorageRepo) DeleteExpiredPastes(keys []string) error {
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

	output, err := storage.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(storage.bucketName),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		log.Printf("Failed to delete objects from bucket %s. Error: %v", storage.bucketName, err)
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
