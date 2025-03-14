package testutils

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

// GetTestData returns sample test data for the metadata table.
func GetTestData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"key":             "test_key",
			"created_at":      time.Now().Add(-1 * time.Hour),
			"expiration_date": time.Now().Add(-1 * time.Minute), // Expired
		},
	}
}

func GetTestEnv() error {
	// Load environment variables from the `.env` file
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}
	return nil
}

var S3TestData = struct {
	Bucket  string
	Region  string
	Key     string
	Content string
}{
	Bucket:  "test-textnest-bucket",
	Region:  "us-east-1",
	Key:     "test_key", // Same key as in metadata
	Content: "Test data for expired paste content",
}
