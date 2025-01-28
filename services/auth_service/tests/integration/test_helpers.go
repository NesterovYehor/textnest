package integration

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func PreparePostgres(ctx context.Context) (*container.PostgresContainer, *sql.DB, error) {
	// Start the Postgres container
	container, err := container.StartPostgres(ctx)
	if err != nil {
		return nil, nil, err
	}

	tableSchema := `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TABLE IF NOT EXISTS users (
            id bigserial PRIMARY KEY,
            created_at timestamp(0) with time zone NOT NULL DEFAULT NOW (),
            name text NOT NULL,
            email citext UNIQUE NOT NULL,
            password_hash bytea NOT NULL,
            activated bool NOT NULL,
            version integer NOT NULL DEFAULT 1
        );
    `

	// Execute SQL to set up the database
	// Get the connection string
	dbUrl, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Println("Error getting connection string:", err)
		return nil, nil, err
	}
	dbUrl += "sslmode=disable"

	// Open the database connection
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("Error opening DB connection:", err)
		return nil, nil, err
	}
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", tableSchema})
	if err != nil {
		fmt.Println("Error executing table creation:", err)
		return nil, nil, err
	}

	// Verify that the table exists
	var exists bool
	err = conn.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	if err != nil {
		fmt.Println("Error checking if table exists:", err)
		return nil, nil, err
	}
	fmt.Println("Users table exists:", exists)

	return container, conn, nil
}
