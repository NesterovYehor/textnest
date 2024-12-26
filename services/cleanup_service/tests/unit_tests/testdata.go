package testutils

import (
	"time"

	"github.com/joho/godotenv"
)

// GetTestData returns sample test data for the metadata table.
func GetTestData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"key":             "expired_key_1",
			"created_at":      time.Now().Add(-1 * time.Hour),
			"expiration_date": time.Now().Add(-1 * time.Minute), // Expired
		},
	}
}

func GetTestEnv() error {
	return godotenv.Load("./.env")
}

var S3TestData = struct {
	Bucket  string
	Region  string
	Key     string
	Content string
}{
	Bucket:  "test-bucket",
	Region:  "us-east-1",
	Key:     "test_key", // Same key as in metadata
	Content: "Test data for expired paste content",
}
