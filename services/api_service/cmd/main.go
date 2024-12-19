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

	http_server "github.com/NesterovYehor/TextNest/pkg/http"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/handler"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	logger, err := setupLogger("app.log")
	if err != nil {
		log.Panic(err)
		return
	}

	cfg, err := config.LoadConfig(ctx, logger)
	if err != nil {
		logger.PrintFatal(ctx, err, nil)
		return
	}

	appContext, err := app.GetAppContext(cfg, ctx, logger)
	if err != nil {
		log.Panic(err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		handler.UploadPaste(w, r, cfg, ctx, appContext) // Pass the pointer to cfg
	})
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
