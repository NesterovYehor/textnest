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
		log.Fatal("Error loading .env.test file")
	}
}

func TestS3Repository_Integration(t *testing.T) {
	// Setup AWS session and S3 client
	region := "us-east-1"
	bucketName := "textnestbucket"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	assert.NoError(t, err)

	s3Client := s3.New(sess)

	// Initialize repository
	repo, err := repository.NewContentRepository(bucketName)

	// Test data
	key := "test_key"
	data := []byte("test_content")

	// Call UploadPasteContent
	err = repo.UploadPasteContent(context.Background(), key, data)
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
