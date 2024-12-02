package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

var timeout = time.Second * 20

type KeyGeneratorRepository struct {
	client  *redis.Client
	breaker *middleware.CircuitBreakerMiddleware
}

func NewRepository(client *redis.Client) *KeyGeneratorRepository {
	cbSettings := gobreaker.Settings{
		Name:        "MetadataRepo",
		MaxRequests: 5,                // Max requests allowed in half-open state
		Interval:    5 * time.Second,  // Time window for tracking errors
		Timeout:     30 * time.Second, // Time to reset the circuit after tripping
	}
	return &KeyGeneratorRepository{
		client:  client,
		breaker: middleware.NewCircuitBreakerMiddleware(cbSettings),
	}
}

// GetKey retrieves a key from the unused_keys set and moves it to used_keys.
func (r *KeyGeneratorRepository) GetKey(ctx context.Context) (string, error) {
	var key string

	operation := func(ctx context.Context) (any, error) {
		// Fetch a key from the unused_keys set
		k, err := r.client.SRandMember(ctx, "unused_keys").Result()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch key from unused_keys: %w", err)
		}
		key = k

		// Transactional move: Remove from unused_keys, add to used_keys
		_, err = r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.SRem(ctx, "unused_keys", key)
			pipe.SAdd(ctx, "used_keys", key)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to move key to used_keys: %w", err)
		}
		return nil, nil
	}

	// Execute the operation with the circuit breaker
	if _, err := r.breaker.Execute(ctx, operation); err != nil {
		return "", err
	}

	// Generate a new key asynchronously
	errCh := make(chan error)
	go func() {
		errCh <- r.generateKey()
		close(errCh)
	}()
	if err := <-errCh; err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}

// ReallocateKey moves a key back to the unused_keys set.
func (r *KeyGeneratorRepository) ReallocateKey(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	operation := func(ctx context.Context) (any, error) {
		_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.SRem(ctx, "used_keys", key)
			pipe.SAdd(ctx, "unused_keys", key)
			return nil
		})
		return nil, err
	}

	if _, err := r.breaker.Execute(ctx, operation); err != nil {
		return fmt.Errorf("circuit breaker triggered during reallocate: %w", err)
	}
	return nil
}

// generateKey creates and stores a unique key in unused_keys.
func (r *KeyGeneratorRepository) generateKey() error {
	const maxRetries = 10
	retries := 0

	for retries < maxRetries {
		key := make([]byte, 12)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("key generation failed: %w", err)
		}

		encodedKey := base64.URLEncoding.EncodeToString(key)[:8]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		operation := func(ctx context.Context) (any, error) {
			// Check for key uniqueness
			isUnused, err := r.isMemberOfSet(encodedKey, "unused_keys")
			if err != nil || isUnused {
				return nil, err
			}

			isUsed, err := r.isMemberOfSet(encodedKey, "used_keys")
			if err != nil || isUsed {
				return nil, err
			}

			// Add key to unused_keys if unique
			if err := r.client.SAdd(ctx, "unused_keys", encodedKey).Err(); err != nil {
				return nil, fmt.Errorf("failed to store new key: %w", err)
			}
			return nil, nil
		}

		if _, err := r.breaker.Execute(ctx, operation); err == nil {
			return nil // Key successfully generated
		}
		retries++
	}
	return errors.New("max retries reached for key generation")
}

// isMemberOfSet checks if a value exists in a Redis set.
func (r *KeyGeneratorRepository) isMemberOfSet(value, setName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	isMember, err := r.client.SIsMember(ctx, setName, value).Result()
	return isMember, err
}

// FillKeys ensures that unused_keys meets a minimum threshold.
func (r *KeyGeneratorRepository) FillKeys(threshold int64) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := r.client.SCard(ctx, "unused_keys").Result()
	if err != nil {
		log.Printf("Failed to get unused_keys count: %v", err)
		return
	}

	for count < threshold {
		if err := r.generateKey(); err != nil {
			log.Printf("Failed to generate key: %v", err)
			break
		}
		count++
	}
}

func IsKeyValid(v *validator.Validator, key string) {
	v.Check(key != "", "key", "key must be provided")
	v.Check(len(key) == 8, "key", "key must be 8 characters long")
}
