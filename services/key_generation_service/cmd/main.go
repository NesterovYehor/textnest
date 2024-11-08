package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/key_generation_service/internal/grpc_server"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()

	redisClient, err := redis.StartRedis(cfg)
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error pinging Redis: %v", err)
	}

	// Check the response and log accordingly
	if pong == "PONG" {
		log.Printf("Successfully connected to Redis: %s", pong)
	} else {
		log.Fatalf("Unexpected Redis response: %s", pong)
	}
	if err != nil {
		log.Panic(err)
	}

	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	key_manager.RegisterKeyManagerServiceServer(grpcSrv.Grpc, key_manager.NewKeyManagerServer(redisClient))

	if err = grpcSrv.RunGrpcServer(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
