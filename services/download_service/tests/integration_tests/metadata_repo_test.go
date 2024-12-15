package integrationtests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

const testKey = "test-key"

var testCreatedAt time.Time

func TestDownloadPasteMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Start Postgres container
	pgContainer, err := container.StartPostgres(ctx)
	assert.Nil(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate Postgres container: %v", err)
		}
	}()

	// Get connection string
	dbUrl, err := pgContainer.ConnectionString(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, dbUrl)
	dbUrl += "sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", dbUrl)
	assert.Nil(t, err)
	defer db.Close()

	// Initialize repository
	repo := repository.NewMetadataRepo(db)
	testCreatedAt := time.Now()
	_, _, err = pgContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", `
        CREATE TABLE IF NOT EXISTS metadata (
            key VARCHAR NOT NULL UNIQUE,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
        );
	`})

	// Insert test data
	_, err = db.ExecContext(ctx, `INSERT INTO metadata(key, created_at, expiration_date) VALUES ($1, $2, $3)`,
		testKey, testCreatedAt, testCreatedAt.Add(time.Minute))
	assert.Nil(t, err)

	// Test DownloadPasteMetadata
	res, err := repo.DownloadPasteMetadata(testKey)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	// Validate returned metadata
	assert.Equal(t, testKey, res.Key)
	assert.Equal(t, testCreatedAt.UTC(), res.CreatedAt.UTC())
	assert.Equal(t, testCreatedAt.Add(time.Minute).UTC(), res.ExpiredDate.UTC())
}

func TestDownloadPasteMetadata_NotFound(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Start Postgres container
	pgContainer, err := container.StartPostgres(ctx)
	assert.Nil(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate Postgres container: %v", err)
		}
	}()
	_, _, err = pgContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", `
        CREATE TABLE IF NOT EXISTS metadata (
            key VARCHAR NOT NULL UNIQUE,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
        );
	`})

	// Get connection string
	dbUrl, err := pgContainer.ConnectionString(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, dbUrl)
	dbUrl += "sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", dbUrl)
	assert.Nil(t, err)
	defer db.Close()

	// Initialize repository
	repo := repository.NewMetadataRepo(db)

	// Test non-existent key
	nonExistentKey := "non-existent-key"
	res, err := repo.DownloadPasteMetadata(nonExistentKey)

	// Assert the expected error and nil result
	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), fmt.Sprintf("no paste found with the key: %s", nonExistentKey))
}
