package testutils

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

// SetupTestDatabase sets up a PostgreSQL container, creates the required table,
// seeds the database with initial test data, and returns the database connection
// and a cleanup function.
func SetupTestDatabase(t assert.TestingT, ctx context.Context) (*sql.DB, func()) {
	// Start PostgreSQL container
	postgresContainer, err := container.StartPostgres(ctx)
	assert.NoError(t, err)

	// Define table schema
	tableSchema := `
        CREATE TABLE IF NOT EXISTS metadata (
            key VARCHAR NOT NULL UNIQUE,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
        );
    `

	// Get the database connection string
	dbURL, err := postgresContainer.ConnectionString(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, dbURL)
	dbURL += "sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	assert.NoError(t, err)

	// Create the table using the provided schema
	_, _, err = postgresContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", tableSchema})
	assert.NoError(t, err)

	// Seed the database with test data
	for _, row := range GetTestData() {
		_, err := db.ExecContext(ctx, `INSERT INTO metadata (key, created_at, expiration_date) VALUES ($1, $2, $3)`,
			row["key"], row["created_at"], row["expiration_date"])
		assert.NoError(t, err, fmt.Sprintf("Failed to insert test data for key: %s", row["key"]))
	}

	// Cleanup function to terminate container and close the database connection
	cleanup := func() {
		defer db.Close()
		defer postgresContainer.Terminate(ctx)
	}

	return db, cleanup
}
