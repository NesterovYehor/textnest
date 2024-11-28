package cache

import (
	"context"
	"fmt"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type redisCache struct {
	client  *redis.Client
	breaker *middleware.CircuitBreakerMiddleware
	ctx     context.Context
}

// NewRedisCache initializes a new Redis cache instance
func NewRedisCache(ctx context.Context, redisAddr string) Cache {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    10 * time.Second, // Time window for tracking errors
		Timeout:     30 * time.Second, // Time to reset the circuit after tripping
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &redisCache{
		client:  rdb,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
		ctx:     ctx,
	}
}

// Set stores Protobuf-serialized data in Redis
func (r *redisCache) Set(key string, value []byte, expiration time.Duration) error {
	operation := func(ctx context.Context) (any, error) {
		// Use the Redis client to set the value with expiration
		err := r.client.Set(ctx, key, value, expiration).Err()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Execute the operation using the Circuit Breaker middleware
	_, err := r.breaker.Execute(r.ctx, operation)
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %w", err)
	}

	return nil
}

func (r *redisCache) Get(key string) ([]byte, bool, error) {
	// Define the operation for the Circuit Breaker
	operation := func(ctx context.Context) (any, error) {
		val, err := r.client.Get(ctx, key).Bytes()
		if err == redis.Nil {
			// Key not found
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		return val, nil
	}

	// Execute the operation with the Circuit Breaker middleware
	result, err := r.breaker.Execute(r.ctx, operation)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get key from Redis: %w", err)
	}

	// If the operation was successful, handle the returned value
	if result == nil {
		return nil, false, nil // Key not found
	}

	return result.([]byte), true, nil
}

func (r *redisCache) Delete(key string) error {
	operation := func(ctx context.Context) (any, error) {
		return nil, r.client.Del(ctx, key).Err()
	}

	_, err := r.breaker.Execute(r.ctx, operation)
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}

	return nil
}

// Clear removes all keys from the current Redis database
func (r *redisCache) Clear() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close closes the Redis client connection
func (r *redisCache) Close() error {
	return r.client.Close()
}
