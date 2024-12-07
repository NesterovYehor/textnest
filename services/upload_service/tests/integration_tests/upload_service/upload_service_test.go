package upload_service_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/logger"
	con "github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUploadService_Upload(t *testing.T) {
	// Set up the test logger and configuration
	log := jsonlog.New(io.Discard, slog.LevelInfo)
	cfg, err := con.LoadConfig(log, context.Background())
	assert.NoError(t, err)

	db, err := sql.Open("postgres", cfg.DBURL)
	assert.NoError(t, err)

	// Initialize real repositories
	metadataRepo := repository.NewMetadataRepository(db)
	storageRepo, err := repository.NewS3Repository(cfg.S3Region)
	assert.NoError(t, err)

	// Start gRPC server
	service := services.NewUploadService(storageRepo, metadataRepo, log, cfg)
	server := grpc.NewServer()
	pb.RegisterUploadServiceServer(server, service)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Grpc.Port)) // Use an unused port
	assert.NoError(t, err)
	go func() {
		err = server.Serve(listener)
		assert.NoError(t, err)
	}()
	defer server.Stop()

	// Create gRPC client
	conn, err := grpc.NewClient(fmt.Sprintf(":%v", cfg.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewUploadServiceClient(conn)

	// Test Upload method
	req := &pb.UploadRequest{
		Key:            "test-key",
		ExpirationDate: timestamppb.New(time.Now().Add(time.Hour)),
		Data:           []byte("test content"),
	}
	resp, err := client.Upload(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Uploaded new paste successfully", resp.Message)
}
