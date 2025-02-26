package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func PreparePostgres(ctx context.Context, tableName, tableSchema string, t *testing.T) (*sql.DB, func()) {
	t.Helper()
	container, err := container.StartPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to start postgres test container: %v", err)
	}

	dbUrl, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fatalf("Failed getting connection string: %v", err)
	}
	dbUrl += "sslmode=disable"

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("Failed opening DB connection: %v", err)
	}
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", tableSchema})
	if err != nil {
		t.Fatalf("Failed executing table creationr %v", err)
	}

	var exists bool
	err = conn.QueryRow(fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%v');", tableName)).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed checking if table exists: %v", err)
	}

	return conn, func() {
		container.Terminate(ctx)
		conn.Close()
	}
}
