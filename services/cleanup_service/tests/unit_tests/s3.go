package testutils

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Returns an error if setup fails, along with a cleanup function to clean up resources.
func SetUpTestS3(ctx context.Context) (func(), error) {
	// Load environment variables from the env.test file
	err := GetTestEnv()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(S3TestData.Region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	// Create the bucket
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(S3TestData.Bucket),
	})
	if err != nil {
		return nil, err
	}

	// Upload test data to the bucket
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(S3TestData.Bucket),
		Key:    aws.String(S3TestData.Key),
		Body:   strings.NewReader(S3TestData.Content),
	})
	if err != nil {
		return nil, err
	}

	// Return cleanup function
	cleanup := func() {
		log.Println("Cleaning up test S3 bucket...")
		err := clearTestS3(ctx)
		if err != nil {
			log.Printf("Failed to clean up test S3 bucket: %v", err)
		}
	}
	return cleanup, nil
}

// clearTestS3 deletes all objects in the test bucket and then deletes the bucket.
func clearTestS3(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(S3TestData.Region))
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg)

	// List all objects in the bucket
	output, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(S3TestData.Bucket),
	})
	if err != nil {
		return err
	}

	// Delete all objects
	for _, obj := range output.Contents {
		_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(S3TestData.Bucket),
			Key:    obj.Key,
		})
		if err != nil {
			log.Printf("Failed to delete object %s: %v", *obj.Key, err)
		}
	}

	// Delete the bucket
	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(S3TestData.Bucket),
	})
	if err != nil {
		return err
	}

	log.Println("Test S3 bucket cleaned up successfully.")
	return nil
}
