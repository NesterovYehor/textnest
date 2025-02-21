package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/handlers"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/scheduler"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
	_ "github.com/lib/pq"
)

// App struct contains all services and components for the application
type App struct {
	Logger              *jsonlog.Logger
	DB                  *sql.DB
	KafkaProducer       *kafka.KafkaProducer
	Scheduler           *scheduler.Checker
	ExpiredPasteHandler *handlers.ExpiredPasteHandler

	Config *config.Config
}

var (
	instance  *App
	once      sync.Once
	initError error // Variable to track initialization errors
)

// NewApp initializes and returns a singleton instance of the application
func NewApp(ctx context.Context) (*App, func(), error) {
	var cleanup func()

	once.Do(func() {
		cfg, err := config.LoadConfig(ctx)
		if err != nil {
			initError = err // Store the error
			cleanup = func() {}
			return
		}

		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			initError = err // Store the error
			cleanup = func() {}
			return
		}
		multiWriter := io.MultiWriter(logFile, os.Stdout)
		logger := jsonlog.New(multiWriter, slog.LevelInfo)

		db, err := openDB(cfg.DBUrl)
		if err != nil {
			initError = err // Store the error
			logFile.Close()
			cleanup = func() {}
			return
		}
		logger.PrintInfo(ctx, "Connected to the database successfully", nil)

		kafkaProducer, err := kafka.NewProducer(*cfg.Kafka, ctx)
		if err != nil {
			initError = err // Store the error
			logFile.Close()
			db.Close()
			cleanup = func() {}
			return
		}
		metadataRepo := repository.NewMetadataRepo(db)
		storageRepo, err := repository.NewStorageRepo(cfg.S3Region, cfg.BucketName)
		if err != nil {
			logger.PrintFatal(ctx, err, nil)
			return
		}

		pasteService := services.NewPasteService(metadataRepo, storageRepo)
		expiredPasteHandler := handlers.NewExpiredPasteHandler(pasteService)

		expirationService := services.NewExpirationService(
			metadataRepo,
			storageRepo,
			kafkaProducer,
		)

		scheduler := scheduler.NewChecker(expirationService, logger)

		instance = &App{
			Logger:              logger,
			DB:                  db,
			KafkaProducer:       kafkaProducer,
			Scheduler:           scheduler,
			ExpiredPasteHandler: expiredPasteHandler,
			Config:              cfg,
		}

		initError = nil // No error occurred
		cleanup = func() {
			logFile.Close()
			db.Close()
			kafkaProducer.Close()
		}
	})

	return instance, cleanup, initError
}

// openDB initializes a new database connection and checks for errors
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (app *App) RunKafkaConsumer(cfg *config.Config, ctx context.Context) error {
	// Map Kafka topics to handlers
	handlers := map[string]kafka.MessageHandler{
		"expired-paste-topic": func(msg *sarama.ConsumerMessage) error {
			// Using the handler for the expired paste topic
			return app.ExpiredPasteHandler.Handle(ctx, msg)
		},
	}

	// Initialize Kafka consumer
	consumer, err := kafka.NewKafkaConsumer(cfg.Kafka, handlers, ctx)
	if err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("Failed to create a new Kafka consumer: %w", err), nil)
		return err
	}

	// Start the consumer
	if err := consumer.Start(); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("Kafka consumer stopped with error: %w", err), nil)
		consumer.Close()
		return err
	}

	return nil
}
