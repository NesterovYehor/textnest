package upload_service_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	postgres "github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	log := jsonlog.New(io.Discard, slog.LevelInfo)
	cfg, err := config.LoadConfig(log, context.Background())
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, err := postgres.Start(ctx, t)
	assert.Nil(t, err)
	_, _, err = postgresContainer.Exec(ctx, []string{"psql", "-U", "testcontainer", "-d", "test_db", "-c", "CREATE TABLE IF NOT EXISTS metadata(key VARCHAR NOT NULL UNIQUE, created_at TIMESTAMP WITH TIME ZONE NOT NULL, expiration_date TIMESTAMP WITH TIME ZONE NOT NULL);"})
	assert.Nil(t, err)
	assert.NotNil(t, postgresContainer)
	dbURL, err := postgresContainer.ConnectionString(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, dbURL)
	dbURL = dbURL + "ssl=disable"

	db, err := sql.Open("postgres", dbURL)
	assert.NoError(t, err)

	// Initialize real repositories
	metadataRepo := repository.NewMetadataRepository(db)
	assert.NotNil(t, metadataRepo)
	storageRepo, err := repository.NewS3Repository(cfg.S3Region)
	assert.NotNil(t, storageRepo)
	assert.NoError(t, err)

	service := services.NewUploadService(storageRepo, metadataRepo, log, cfg)

	resp, err := service.Upload(context.Background(), &grpcRequest)
	assert.Nil(t, err)
	assert.Equal(t, &grpcResponse, resp)
}
