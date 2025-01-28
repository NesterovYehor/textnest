package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/textnest/services/auth_service/api"
	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/controlers"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := setupLogger("app.log")
	if err != nil {
		fmt.Println(err)
	}

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.PrintFatal(ctx, err, nil)
		return
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.PrintFatal(ctx, err, nil)
		return
	}

	userModel := models.NewUserModel(db)
	userSrv := services.NewUserService(userModel)
	tokenSrv := services.NewJwtService(cfg.JwtConfig)
	controler := controlers.NewAuthControler(logger, userSrv, tokenSrv)

	server := grpc.NewGrpcServer(cfg.Grpc)
	pb.RegisterAuthServiceServer(server.Grpc, controler)
	if err := server.RunGrpcServer(ctx); err != nil {
		logger.PrintError(ctx, err, nil)
		return
	}
}

// setupLogger initializes the application logger
func setupLogger(logFilePath string) (*jsonlog.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return nil, err
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	return jsonlog.New(multiWriter, slog.LevelInfo), nil
}
