package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/pkg/kafka"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/redis"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/services"
	key_manager "github.com/NesterovYehor/TextNest/services/key_generation_service/proto"
)

func main() {
	// Setup signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize logger
	log, err := setupLogger("app.log")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(ctx, log)
	if err != nil {
		log.PrintError(ctx, fmt.Errorf("failed to load configuration: %v", err), nil)
		return
	}

	// Initialize Redis client
	redisClient, err := redis.StartRedis(cfg)
	if err != nil {
		log.PrintError(ctx, err, nil)
		return
	}

	// Initialize key management repository
	repo := repository.NewRepository(redisClient)
	repo.FillKeys(10) // Ensure a minimum threshold of unused keys

	// Start gRPC server
	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)
	keyManagerService := services.NewKeyManagerServer(repo)
	key_manager.RegisterKeyGeneratorServer(grpcSrv.Grpc, keyManagerService)

	// Start Kafka consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := startKafkaConsumer(log, cfg, ctx, repo)
		if err != nil {
			log.PrintError(ctx, fmt.Errorf("Kafka consumer error: %v", err), nil)
		}
	}()

	// Start gRPC server in a separate goroutine
	go func() {
		if err := grpcSrv.RunGrpcServer(ctx); err != nil {
			log.PrintError(ctx, fmt.Errorf("gRPC server error: %v", err), nil)
		}
	}()

	// Wait for all services to finish
	wg.Wait()
	log.PrintInfo(ctx, "All services have shut down gracefully.", nil)
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

// startKafkaConsumer initializes and starts the Kafka consumer
// main.go

// startKafkaConsumer initializes and starts the Kafka consumer
func startKafkaConsumer(log *jsonlog.Logger, cfg *config.Config, ctx context.Context, repo *repository.KeyGeneratorRepository) error {
	handlers := map[string]kafka.MessageHandler{
		"delete-expired-paste-topic": func(msg *sarama.ConsumerMessage) error {
			// Additional logging for context
			log.PrintInfo(ctx, fmt.Sprintf("Consumed message from topic %s: %s", msg.Topic, string(msg.Value)), nil)
			return repo.ReallocateKey(string(msg.Value))
		},
	}

	consumer, err := kafka.NewKafkaConsumer(&cfg.Kafka, handlers, ctx)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	// Ensure graceful shutdown by capturing context
	go func() {
		<-ctx.Done()
		log.PrintInfo(ctx, "Shutting down Kafka consumer...", nil)
		consumer.Close() // Ensure proper closing of the consumer
	}()

	if err := consumer.Start(); err != nil {
		consumer.Close() // Make sure the consumer is closed if an error occurs
		return fmt.Errorf("Kafka consumer stopped with error: %w", err)
	}

	log.PrintInfo(ctx, "Kafka consumer started", nil)
	return nil
}
