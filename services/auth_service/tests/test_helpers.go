package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func PreparePostgres(ctx context.Context, tableName, tableSchema string, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	container, err := container.StartPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to start postgres test container: %v", err)
	}

	dbLink, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Println("Error ", err)
		t.Fatalf("Failed getting connection string: %v", err)
	}
	dbLink += "sslmode=disable"

	conn, err := pgxpool.New(ctx, dbLink)
	if err != nil {
		t.Fatalf("Failed opening DB connection: %v", err)
	}
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", tableSchema})
	if err != nil {
		t.Fatalf("Failed executing table creationr %v", err)
	}

	var exists bool
	err = conn.QueryRow(ctx, fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%v');", tableName)).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed checking if table exists: %v", err)
	}

	return conn, func() {
		container.Terminate(ctx)
		conn.Close()
	}
}
