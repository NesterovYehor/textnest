package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func PreparePostgres(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	t.Helper()
	// Start the Postgres container
	container, err := container.StartPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to start postgres test container: %v", err)
	}

	tableSchema := `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
            created_at timestamp(0) with time zone NOT NULL DEFAULT NOW (),
            name text NOT NULL,
            email citext UNIQUE NOT NULL,
            password_hash bytea NOT NULL,
            activated bool NOT NULL DEFAULT false
        );
    `

	// Execute SQL to set up the database
	// Get the connection string
	dbUrl, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fatalf("Failed getting connection string: %v", err)
	}
	dbUrl += "sslmode=disable"

	// Open the database connection
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("Failed opening DB connection: %v", err)
	}
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", tableSchema})
	if err != nil {
		t.Fatalf("Failed executing table creationr %v", err)
	}

	// Verify that the table exists
	var exists bool
	err = conn.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed checking if table exists: %v", err)
	}

	return conn, func() {
		container.Terminate(ctx)
		conn.Close()
	}
}
