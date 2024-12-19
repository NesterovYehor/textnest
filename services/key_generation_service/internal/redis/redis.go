package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	"github.com/redis/go-redis/v9"
)

func StartRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Printf("redis server is on: %s\n", cfg.RedisAddr)

	return rdb, nil
}
