package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	testutils "github.com/NesterovYehor/TextNest/services/cleanup_service/tests/unit_tests"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAndReturnExpiredKeys(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	db, cleanup := testutils.SetupTestDatabase(t, ctx)
	defer cleanup()

	// Create the repository and call DeleteAndReturnExpiredKeys
	factory := repository.NewRepositoryFactory(db)
	repo := factory.CreateMetadataRepository()
	expiredKeys, err := repo.DeleteAndReturnExpiredKeys()
	assert.NoError(t, err)
	assert.Equal(t, []string{"test_key"}, expiredKeys, "Expected expired keys to match the inserted key")

	// Validate that the expired row has been deleted
	exists := testutils.VerifyRowExists(t, db, "test_key")
	assert.False(t, exists, "Expected the expired key to be deleted from the database")
}

func TestDeletePasteByKey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	db, cleanup := testutils.SetupTestDatabase(t, ctx)
	defer cleanup()

	// Create the repository and call DeletePasteByKey
	factory := repository.NewRepositoryFactory(db)
	repo := factory.CreateMetadataRepository()
	err := repo.DeletePasteByKey("test_key")
	assert.NoError(t, err)

	// Validate that the key has been deleted
	exists := testutils.VerifyRowExists(t, db, "test_key")
	assert.False(t, exists, "Expected the key to be deleted from the database")
}
