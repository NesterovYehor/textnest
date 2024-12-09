package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUploadPasteMetadata(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Start a PostgreSQL container
	postgresContainer, err := Start(ctx, t)
	assert.Nil(t, err)

	// Create metadata table
	_, _, err = postgresContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", `
		CREATE TABLE IF NOT EXISTS metadata (
			key VARCHAR NOT NULL UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`})
	assert.Nil(t, err)

	// Get the database connection
	dbURL, err := postgresContainer.ConnectionString(ctx)
	assert.Nil(t, err)
	dbURL += "sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	assert.Nil(t, err)

	// Initialize repository
	metadataRepo := repository.NewMetadataRepository(db)
	assert.NotNil(t, metadataRepo)

	// Define metadata to insert
	metadata := &models.MetaData{
		Key:            "test_key",
		CreatedAt:      time.Now(),
		ExpirationDate: time.Now().Add(24 * time.Hour),
	}

	// Test the UploadPasteMetadata method
	err = metadataRepo.UploadPasteMetadata(ctx, metadata)
	assert.Nil(t, err)
}

type PostgresContainer struct {
	postgres.PostgresContainer
}

func Start(ctx context.Context, t *testing.T) (*PostgresContainer, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("testcontainer"),
		postgres.WithPassword("testcontainer"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		return nil, err
	}
	return &PostgresContainer{*postgresContainer}, nil
}
