package integrationtests

import (
	"log"
	"strings"
	"testing"

	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	key        = "test-key"
	region     = "us-east-1"
	bucketName = "test-textnest-bucket"
	content    = "test-content"
)

func init() {
	// Load environment variables from the env.test file
	err := godotenv.Load("../../.env.test") // Adjust the path as needed
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestStorageRepository(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	assert.NoError(t, err)

	s3Client := s3.New(sess)

	// Create the S3 bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	assert.NoError(t, err) // Ensure no errors during `PutObject`

	// Delete bucket after objects
	defer func() {
		resp, err := s3Client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err) // Add error handling for `ListObjects`

		for _, obj := range resp.Contents {
			_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*obj.Key),
			})
			assert.NoError(t, err) // Add error handling for `DeleteObject`
		}

		_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err) // Add error handling for `DeleteBucket`
	}()

	repo, err := repository.NewStorageRepository(bucketName)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	data, err := repo.DownloadPasteContent(bucketName, key)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestStorageRepository_NoFound(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	assert.NoError(t, err)

	s3Client := s3.New(sess)

	// Create the S3 bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	// Delete bucket after objects
	defer func() {
		resp, err := s3Client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err) // Add error handling for `ListObjects`

		for _, obj := range resp.Contents {
			_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*obj.Key),
			})
			assert.NoError(t, err) // Add error handling for `DeleteObject`
		}

		_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		assert.NoError(t, err) // Add error handling for `DeleteBucket`
	}()

	repo, err := repository.NewStorageRepository(bucketName)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	_, err = repo.DownloadPasteContent(bucketName, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker triggered")
	assert.Contains(t, err.Error(), "NoSuchKey")
}
