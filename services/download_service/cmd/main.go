package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	download_service "github.com/NesterovYehor/TextNest/services/download_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/repository"
	_ "github.com/lib/pq" // PostgreSQL driver
)

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
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	log := jsonlog.New(logFile, slog.LevelInfo)

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

	// Register the UploadService with the gRPC server
	download_service.RegisterDownloadServiceServer(grpcSrv.Grpc, download_service.NewDownloadService(storageRepo, metadataRepo, log, cfg, ctx))

	// Run the gRPC server
	if err := grpcSrv.RunGrpcServer(ctx); err != nil {
		log.PrintFatal(ctx, err, nil)
		return
	}
}
