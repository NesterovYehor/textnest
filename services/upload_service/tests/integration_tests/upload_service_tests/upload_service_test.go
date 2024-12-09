package upload_service_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresContainer struct {
	postgres.PostgresContainer
}

var grpcRequest = pb.UploadRequest{
	Key:            "test_key",
	ExpirationDate: timestamppb.New(time.Now().Add(time.Minute)),
	Data:           []byte("test_content"),
}

var grpcResponse = &pb.UploadResponse{
	Message: "Uploaded new paste successfully",
}

func TestUploadService_Upload(t *testing.T) {
	// Set up the test logger and configuration
	log := jsonlog.New(os.Stdout, slog.LevelInfo) // Changed to os.Stdout for visible logs
	cfg, err := config.LoadConfig(log, context.Background())
	if err != nil {
		log.PrintFatal(context.Background(), fmt.Errorf("failed to load config: %w", err), nil)
		t.FailNow()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Postgres container
	postgresContainer, err := Start(ctx, t)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to start postgres container: %w", err), nil)
		t.FailNow()
	}
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.PrintError(ctx, fmt.Errorf("failed to terminate postgres container: %w", err), nil)
		}
	}()

	// Create the metadata table
	_, _, err = postgresContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", `
		CREATE TABLE IF NOT EXISTS metadata (
			key VARCHAR NOT NULL UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			expiration_date TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`})
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to create metadata table: %w", err), nil)
		t.FailNow()
	}

	// Connect to the database
	dbURL, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to get connection string: %w", err), nil)
		t.FailNow()
	}
	dbURL += "sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to connect to the database: %w", err), nil)
		t.FailNow()
	}

	// Initialize real repositories
	metadataRepo := repository.NewMetadataRepository(db)
	if metadataRepo == nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to initialize metadata repository"), nil)
		t.FailNow()
	}
	storageRepo, err := repository.NewS3Repository(cfg.S3Region)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to initialize S3 repository: %w", err), nil)
		t.FailNow()
	}

	// Initialize the service and make a gRPC request
	service := services.NewUploadService(storageRepo, metadataRepo, log, cfg)
	resp, err := service.Upload(context.Background(), &grpcRequest)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to execute Upload method: %w", err), nil)
		t.FailNow()
	}
	assert.Equal(t, grpcResponse, resp)
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
		t.Logf("failed to start postgres container: %v", err)
		return nil, err
	}
	return &PostgresContainer{*postgresContainer}, nil
}
