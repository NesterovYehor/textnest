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
	log "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/download_service/internal/coordinators"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	log, err := setupLogger("app.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(log, ctx)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return
	}

	db, err := initializeDatabase(cfg.DBURL, log, ctx)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return
	}

	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	coord, err := coordinators.NewDownloadCoordinator(ctx, cfg, log, db)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return
	}
	pb.RegisterPasteDownloadServer(grpcSrv.Grpc, coord)

	if err := grpcSrv.RunGrpcServer(ctx); err != nil {
		log.PrintFatal(ctx, err, nil)
		return
	}
}

func setupLogger(logFilePath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		logFile.Close()
		return nil, err
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	return log.New(multiWriter, slog.LevelInfo), nil
}

func initializeDatabase(dsn string, log *log.Logger, ctx context.Context) (*sql.DB, error) {
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
