package redis

import (
	"context"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/config"
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

	return rdb, nil
}
