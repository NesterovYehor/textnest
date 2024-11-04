package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/keymanager"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/routes"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Config *config.Config
}

func (server *Server) Start() error {
	client, err := startRedis(server.Config)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	router := http.NewServeMux()
	routes.SetupRoutes(router, client)

	srv := http.Server{
		Addr:         server.Config.Addr,
		Handler:      router,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	fmt.Printf("\nStarting KGS Server on port: %v\n", server.Config.Addr)
	return srv.ListenAndServe()
}

func startRedis(cfg *config.Config) (*redis.Client, error) {
	redisOpts := redis.Options{
		Addr:     cfg.RedisOption.Addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	}
	client := redis.NewClient(&redisOpts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	keymanager.FillKeys(client, 10)

	return client, nil
}
