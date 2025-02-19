package cache

import (
	"context"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
    "google.golang.org/protobuf/proto"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type redisCache struct {
	client     *redis.Client
	breaker    *middleware.CircuitBreakerMiddleware
	expiration time.Duration
}

// NewRedisCache initializes a new Redis cache instance
func NewRedisCache(redisAddr string) (Cache, error) {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,
		Interval:    5 * time.Second,
		Timeout:     30 * time.Second,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       0,
		PoolSize: 10,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &redisCache{
		client:     rdb,
		expiration: time.Hour * 24,
		breaker:    middleware.NewCircuitBreakerMiddleware(cbSettings),
	}, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value *pb.Metadata) error {
	operation := func(ctx context.Context) (any, error) {
		data, err := proto.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal value to bytes:%v", err)
		}
		if err := r.client.Set(ctx, key, data, r.expiration).Err(); err != nil {
			return nil, fmt.Errorf("Failed to store data in cache: %v", err)
		}
		return nil, nil
	}

	_, err := r.breaker.Execute(ctx, operation)
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %w", err)
	}

	return nil
}

func (r *redisCache) Get(ctx context.Context, key string) (*pb.Metadata, bool, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	metadata := &pb.Metadata{}
	if err := proto.Unmarshal(val, metadata); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal data from cache: %v", err)
	}

	return metadata, true, nil
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	operation := func(ctx context.Context) (any, error) {
		return nil, r.client.Del(ctx, key).Err()
	}

	_, err := r.breaker.Execute(ctx, operation)
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}

	return nil
}

func (r *redisCache) Clear(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
