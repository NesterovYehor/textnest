package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/redis"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/routes"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()

	redisClient, err := redis.StartRedis(cfg)
	if err != nil {
		log.Panic(err)
	}

    grpcSrv := grpc.

    if err := grpc.{
		log.Fatalf("Server failed to start: %v", err)
	}
}
