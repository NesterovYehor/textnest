package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	upload_service "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/models"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/storage"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// openDB initializes a new database connection and checks for errors
func openDB(dsn string) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Setup graceful shutdown on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize gRPC server
	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	// Initialize S3 storage
	storage, err := storage.NewS3Storage(cfg.Storage.Bucket, cfg.Storage.Region)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize the database connection using openDB function
	db, err := openDB(cfg.DbUrl) // Make sure cfg.Database.DSN contains your correct DSN
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	// Initialize models with the database connection
	models := models.NewModel(db)

	// Register the UploadService with the gRPC server
	upload_service.RegisterUploadServiceServer(grpcSrv.Grpc, upload_service.NewUploadServer(storage, models))

	// Run the gRPC server
	if err := grpcSrv.RunGrpcServer(ctx); err != nil {
		log.Fatal(err)
		return
	}
}
