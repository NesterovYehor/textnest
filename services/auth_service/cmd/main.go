package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	pb "github.com/NesterovYehor/textnest/services/auth_service/api"
	"github.com/NesterovYehor/textnest/services/auth_service/config"
	controllers "github.com/NesterovYehor/textnest/services/auth_service/internal/controlers"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/database"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/mailer"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
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

	db, err := database.New(cfg.DB, ctx)
	if err != nil {
		logger.PrintFatal(ctx, err, nil)
		return
	}
	defer db.Close()
	model := models.New(db.Pool)

	userSrv := services.NewUserService(model.User)
	tokenSrv := services.NewTokenService(cfg.JwtConfig, model.Token)
	mailer := mailer.NewMailer(cfg)
	controler := controllers.NewAuthController(logger, userSrv, tokenSrv, mailer)

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
