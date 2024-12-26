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

	appContext, err := app.GetAppContext(cfg, ctx, logger)
	if err != nil {
		logger.PrintFatal(ctx, fmt.Errorf("failed to initialize app context: %w", err), nil)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/upload", uploadPasteHandler(cfg, appContext, logger))
	mux.HandleFunc("/v1/download", downloadPasteHandler(cfg, appContext, logger))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: mux,
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

func downloadPasteHandler(cfg *config.Config, appContext *app.AppContext, logger *jsonlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.PrintInfo(r.Context(), "Request received", nil)
		handler.DownloadPaste(w, r, cfg, r.Context(), appContext)
	}
}

func uploadPasteHandler(cfg *config.Config, appContext *app.AppContext, logger *jsonlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.PrintInfo(r.Context(), "Request received", nil)
		handler.UploadPaste(w, r, cfg, r.Context(), appContext)
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
