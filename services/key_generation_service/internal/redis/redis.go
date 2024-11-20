package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/config"
	"github.com/redis/go-redis/v9"
)

func StartRedis(cfg *config.Config) (*redis.Client, error) {
	redisOpts := redis.Options{
		Addr:     cfg.RedisOption.Addr,
		Password: "",
		DB:       0,
	}
	rdb := redis.NewClient(&redisOpts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Printf("redis server is on: %s\n", cfg.RedisOption.Addr)

	keymanager.FillKeys(rdb, 10)

	return rdb, nil
}
