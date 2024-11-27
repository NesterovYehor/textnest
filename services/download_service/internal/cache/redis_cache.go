package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache initializes a new Redis cache instance
func NewRedisCache(ctx context.Context, redisAddr string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &redisCache{
		client: rdb,
		ctx:    ctx,
	}
}

// Set stores Protobuf-serialized data in Redis
func (r *redisCache) Set(key string, value []byte, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err() // Store the serialized data in Redis with the specified expiration time
}

func (r *redisCache) Get(key string) ([]byte, bool, error) {
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil // Key not found
	} else if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

// Delete removes a key from Redis
func (r *redisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Clear removes all keys from the current Redis database
func (r *redisCache) Clear() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close closes the Redis client connection
func (r *redisCache) Close() error {
	return r.client.Close()
}
