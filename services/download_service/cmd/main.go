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

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	pb "github.com/NesterovYehor/TextNest/services/download_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	log, err := setupLogger("app.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	// Setup graceful shutdown on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize configuration
	cfg, err := config.LoadConfig(log, ctx)
	if err != nil {
		log.PrintError(ctx, err, nil)
	}

	// Initialize gRPC server
	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	// Initialize S3 storage
	storageRepo, err := repository.NewStorageRepository(cfg.BucketName, cfg.S3Region)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return
	}

	// Initialize the database connection using openDB function
	db, err := openDB(cfg.DBURL) // Make sure cfg.Database.DSN contains your correct DSN
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("Failed to connect to the database:", err), nil)
		return
	}
	defer db.Close()

	// Initialize models with the database connection
	metadataRepo := repository.NewMetadataRepo(db)
	newDwonloadServie, err := services.NewDownloadService(storageRepo, metadataRepo, log, cfg, ctx)
	if err != nil {
		log.PrintFatal(ctx, err, nil)
		return
	}

	// Register the UploadService with the gRPC server
	pb.RegisterDownloadServiceServer(grpcSrv.Grpc, newDwonloadServie)

	// Run the gRPC server
	if err := grpcSrv.RunGrpcServer(ctx); err != nil {
		log.PrintFatal(ctx, err, nil)
		return
	}
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
