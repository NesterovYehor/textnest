package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/key_generation_service/internal/grpc_server/protos"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/handler"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()

	redisClient, err := redis.StartRedis(cfg)
	if err != nil {
		log.Panic(err)
	}

	grpcSrv := grpc.NewGrpcServer(cfg.Grpc)

	key_manager.RegisterKeyManagerServiceServer(grpcSrv.Grpc, handler.NewKeyManagerServer(redisClient))

	if err = grpcSrv.RunGrpcServer(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
