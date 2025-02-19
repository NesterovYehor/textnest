package tests

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/kafka"
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

func SetUpRedis(ctx context.Context, t *testing.T) (string, func()) {
	t.Helper()

	container, err := container.StartRedis(ctx, &container.RedisContainerOpts{Addr: 6379})
	if err != nil {
		t.Fatalf("Failed to start redis test container: %v", err)
	}
	conn, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get redis test container connection string: %v", err)
	}
	parsedConn := strings.TrimPrefix(conn, "redis://")

	return parsedConn, func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate redis test container: %v", err)
		}
	}
}

func SetUpKafka(ctx context.Context, t *testing.T) *kafka.KafkaProducer {
	t.Helper()

	kafkaSetup, err := container.StartKafka(ctx)
	if err != nil {
		t.Fatalf("Failed to start Kafka test container: %v", err)
	}

	if len(kafkaSetup.BrokerAddr) == 0 {
		t.Fatal("No Kafka broker address received from test container")
	}

	// Use the correct broker address in Kafka config
	kafkaConf := kafka.LoadKafkaConfig(
		[]string{kafkaSetup.BrokerAddr[0]}, // Pass the correct test broker
		[]string{"test-topic"},
		"test-group",
		5,
	)
	kafkaProd, err := kafka.NewProducer(*kafkaConf, ctx)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}

	return kafkaProd
}
