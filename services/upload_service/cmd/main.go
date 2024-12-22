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
	pb "github.com/NesterovYehor/TextNest/services/upload_service/internal/api"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/upload_service/internal/coordinators"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log, err := setupLogger("app.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	cfg, err := config.LoadConfig(log, ctx)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to load configuration: %w", err), nil)
		return
	}

	db, err := initializeDatabase(cfg.DBURL, log, ctx)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to initialize database: %w", err), nil)
		return
	}
	defer db.Close()

	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)
	coord, err := coordinators.NewUploadCoordinator(cfg, log, db)
	if err != nil {
		log.PrintFatal(ctx, fmt.Errorf("failed to initialize coordinator: %w", err), nil)
		return
	}
	pb.RegisterPasteUploadServer(grpcSrv.Grpc, coord)

	log.PrintInfo(ctx, "Starting gRPC server", nil)

	go func() {
		if err := grpcSrv.RunGrpcServer(ctx); err != nil {
			log.PrintFatal(ctx, fmt.Errorf("gRPC server error: %w", err), nil)
		}
	}()

	<-ctx.Done()
	log.PrintInfo(ctx, "Shutting down service...", nil)
	grpcSrv.Grpc.GracefulStop()
	log.PrintInfo(ctx, "Service stopped gracefully", nil)
}

func setupLogger(logFilePath string) (*jsonlog.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	return jsonlog.New(multiWriter, slog.LevelInfo), nil
}

func initializeDatabase(dsn string, log *jsonlog.Logger, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.PrintInfo(ctx, "Connected to the database", nil)
	return db, nil
}
