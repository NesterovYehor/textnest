package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	upload_service "github.com/NesterovYehor/TextNest/services/upload_service/internal/grpc"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Set up context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Open log file
	log, err := setupLogger("app.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(log, ctx)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to load configuration: %w", err), nil)
		return
	}

	// Initialize database connection
	db, err := initializeDatabase(cfg.DBURL, log, ctx)
	if err != nil {
		return
	}
	defer db.Close()

	// Initialize S3 storage
	storageRepo, err := initializeS3Storage(cfg.BucketName, log, ctx)
	if err != nil {
		log.PrintFatal(ctx, err, nil)
		return
	}

	// Initialize metadata repository
	metadataRepo := repository.NewMetadataRepository(db)

	// Initialize gRPC server
	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)
	uploadService := services.NewUploadService(storageRepo, metadataRepo, log, cfg)
	upload_service.RegisterUploadServiceServer(grpcSrv.Grpc, uploadService)

	log.PrintInfo(ctx, "Starting gRPC server", nil)

	// Run gRPC server
	if err := grpcSrv.RunGrpcServer(ctx); err != nil {
		log.PrintFatal(ctx, fmt.Errorf("gRPC server encountered an error: %w", err), nil)
		return
	}

	log.PrintInfo(ctx, "gRPC server shut down gracefully", nil)
}

// setupLogger initializes the application logger
func setupLogger(logFilePath string) (*jsonlog.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		logFile.Close()
		return nil, err
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	return jsonlog.New(multiWriter, slog.LevelInfo), nil
}

// Initializes and verifies database connection.
func initializeDatabase(dsn string, log *jsonlog.Logger, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.PrintError(ctx, fmt.Errorf("failed to connect to the database: %w", err), nil)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.PrintInfo(ctx, "Connected to the database", nil)
	return db, nil
}

// Initializes S3 storage repository.
func initializeS3Storage(bucketName string, log *jsonlog.Logger, ctx context.Context) (repository.StorageRepository, error) {
	storageRepo, err := repository.NewS3Repository(bucketName)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("failed to initialize S3 storage: %w", err), nil)
		return nil, err
	}
	log.PrintInfo(ctx, "S3 storage initialized successfully", nil)
	return storageRepo, nil
}
