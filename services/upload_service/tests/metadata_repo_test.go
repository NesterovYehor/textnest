package tests

import (
	"context"
	"testing"
	"time"

	pb "github.com/NesterovYehor/TextNest/services/upload_service/api"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var testData = &pb.UploadPasteRequest{
	Key:            "test-key",
	UserId:         "test-userid",
	ExpirationDate: timestamppb.New(time.Now().Add(time.Hour)),
	Title:          "test-title",
}

func TestUploadPaste(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Start database
	db, cleanup := SetUpPostgres(ctx, t)
	defer cleanup()

	// Initialize repository
	repo := repository.NewMetadataRepository(db)

	// Insert metadata
	err := repo.InsertPasteMetadata(ctx, testData)
	assert.NoError(t, err, "Failed to insert paste metadata")

	// Query DB directly to verify insertion
	var key, title, userId string
	err = db.QueryRowContext(ctx, "SELECT key, title, user_id FROM metadata WHERE key = $1", testData.Key).
		Scan(&key, &title, &userId)
	assert.NoError(t, err, "Failed to retrieve inserted metadata")
	assert.Equal(t, testData.Key, key, "Metadata key mismatch")
	assert.Equal(t, testData.Title, title, "Metadata title mismatch")
	assert.Equal(t, testData.UserId, userId, "Metadata user ID mismatch")
}

func TestUpdatePaste(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, cleanup := SetUpPostgres(ctx, t)
	defer cleanup()

	repo := repository.NewMetadataRepository(db)

	// Insert test data
	assert.NoError(t, repo.InsertPasteMetadata(ctx, testData))

	// Perform the update
	newExpiration := time.Now().Add(time.Hour) // Set expiration an hour ahead
	assert.NoError(t, repo.UpdatePasteMetadata(ctx, newExpiration, testData.Key))

	// Verify the update
	var expirationDate time.Time
	err := db.QueryRowContext(ctx, "SELECT expiration_date FROM metadata WHERE key = $1", testData.Key).
		Scan(&expirationDate)
	assert.NoError(t, err, "Failed to retrieve updated metadata")

	// Check if the expiration date matches (allowing a small margin)
	assert.WithinDuration(t, newExpiration, expirationDate, time.Second, "Expiration date mismatch")
}
