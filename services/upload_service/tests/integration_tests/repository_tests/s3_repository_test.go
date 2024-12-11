package repository_test

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Load environment variables from the env.test file
	err := godotenv.Load("../../../.env.test") // Adjust the path as needed
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestS3Repository_Integration(t *testing.T) {
	// Setup AWS session and S3 client
	region := "us-east-1"
	bucketName := "test-textnest-bucket"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	assert.NoError(t, err)

	s3Client := s3.New(sess)

	// Create the S3 bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	assert.NoError(t, err)
	defer func() {
		// Delete all objects in the bucket before deleting the bucket itself
        resp, err := s3Client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			log.Println("Error listing objects:", err)
		} else {
			for _, obj := range resp.Contents {
				_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(*obj.Key),
				})
				if err != nil {
					log.Println("Error deleting object:", err)
				}
			}
		}

		// Now delete the bucket
		_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
	}()
	// Initialize repository
	repo, err := repository.NewS3Repository(region)

	// Test data
	key := "test_key"
	data := []byte("test_content")

	// Call UploadPasteContent
	err = repo.UploadPasteContent(context.Background(), bucketName, key, data)
	assert.NoError(t, err)

	// Verify content in S3
	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	assert.NoError(t, err)

	// Check content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, data, buf.Bytes())
}
