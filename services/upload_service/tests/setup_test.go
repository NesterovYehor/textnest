package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func SetUpPostgres(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	t.Helper() // Marks this function as a test helper

	// Start test container
	pgContainer, err := container.StartPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to start postgres test container: %v", err)
	}

	// Get connection string
	dbURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}
	fmt.Println(dbURL)

	// Open DB connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}

	// Run migrations (or a minimal schema for testing)
	query := `
    CREATE TABLE IF NOT EXISTS metadata (
        key VARCHAR NOT NULL UNIQUE,
        title TEXT,
        user_id TEXT DEFAULT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
        );

    `
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		t.Fatalf("Failed to run test migration: %v", err)
	}

	return db, func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}
}
