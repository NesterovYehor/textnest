package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/key_generation_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/kafka"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/redis"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()

	var wg sync.WaitGroup

	// Initialize Redis
	redisClient, err := redis.StartRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil || pong != "PONG" {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Initialize gRPC Server
	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	// Initialize Repository
	repo := repository.NewRepository(redisClient)

	// Register gRPC Key Manager Service
	key_manager.RegisterKeyManagerServiceServer(grpcSrv.Grpc, key_manager.NewKeyManagerServer(redisClient, repo))

	// Initialize Kafka Consumer
	kafkaConsumer, err := kafka.NewKeyReallocatorConsumer(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.ConsumerTopic, repo)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start Kafka Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer encountered an error: %v", err)
		}
	}()

	// Start gRPC Server
	go func() {
		if err := grpcSrv.RunGrpcServer(ctx); err != nil {
			log.Fatalf("gRPC server encountered an error: %v", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All services have shut down gracefully.")
}
