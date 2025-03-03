// @title TextNest API
// @version 1.0
// @description This is the API documentation for the TextNest project.
// @contact.name Your Name
// @contact.email your.email@example.com
// @host localhost:8080
// @BasePath /v1
// @schemes http https

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/handler"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/middlewares"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := setupLogger("app.log")
	if err != nil {
		log.Panic(err)
		return
	}

	cfg, err := config.LoadConfig(ctx, logger)
	if err != nil {
		logger.PrintFatal(ctx, fmt.Errorf("failed to load config: %w", err), nil)
		return
	}

	appContext, err := app.NewAppContext(cfg, ctx, logger)
	if err != nil {
		logger.PrintFatal(ctx, fmt.Errorf("failed to initialize app context: %w", err), nil)
		return
	}
	defer func() {
		if err := appContext.Close(); err != nil {
			logger.PrintError(ctx, err, nil)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("POST /v1/pastes/upload", middlewares.Authenticate(http.HandlerFunc(handler.UploadPasteHandler(appContext))))
	mux.Handle("GET /v1/pastes/download", middlewares.Authenticate(handler.DownloadPaste(cfg, appContext)))
	mux.Handle("GET /v1/pastes/download/all", middlewares.Authenticate(handler.DownloadAllPastesOfUser(cfg, appContext)))
	mux.Handle("GET /v1/pastes/update/{key}", middlewares.Authenticate(handler.UpdatePasteHandler(appContext)))
	mux.Handle("DELETE /v1/pastes/expire/{key}", middlewares.Authenticate(handler.ExpirePasteHandler(appContext)))
	mux.Handle("DELETE /v1/pastes/expire/all", middlewares.Authenticate(handler.ExpireAllUserPastesHandler(appContext)))
	mux.HandleFunc("POST /v1/users/signup", handler.SignUpHandler(appContext, ctx))
	mux.HandleFunc("GET /v1/users/login", handler.LogInHandler(appContext, ctx))
	mux.HandleFunc("GET /v1/users/activate/{token}", handler.ActivateUser(appContext))
	mux.Handle("POST /v1/users/password/{token}", middlewares.Authenticate(handler.ResetPassword(appContext)))
	mux.Handle("GET /v1/tokens/password-reset", middlewares.Authenticate(handler.SendPasswordResetEmail(appContext)))
	mux.Handle("POST /v1/tokens/refresh", http.HandlerFunc(handler.RefreshTokens(appContext)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         cfg.HttpAddr,
		Handler:      middlewares.RateLimit(cfg, mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.PrintInfo(ctx, fmt.Sprintf("Starting server on %v", cfg.HttpAddr), nil)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.PrintError(ctx, fmt.Errorf("HTTP server error: %w", err), nil)
		}
	}()

	<-ctx.Done()
	logger.PrintInfo(ctx, "Shutting down server", nil)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.PrintError(ctx, fmt.Errorf("error during shutdown: %w", err), nil)
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
