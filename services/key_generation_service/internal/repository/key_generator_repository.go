package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/redis/go-redis/v9"
)

var timeout = time.Second * 20

type KeyGeneratorRepository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *KeyGeneratorRepository {
	return &KeyGeneratorRepository{client: client}
}

// GetKey retrieves a key from the unused_keys set and moves it to used_keys.
func (r *KeyGeneratorRepository) GetKey() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	key, err := r.client.SRandMember(ctx, "unused_keys").Result()
	if err != nil {
		return "", fmt.Errorf("failed to fetch key from unused_keys: %w", err)
	}

	// Transactional move: Remove from unused_keys, add to used_keys
	_, err = r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, "unused_keys", key)
		pipe.SAdd(ctx, "used_keys", key)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to move key to used_keys: %w", err)
	}

	errCh := make(chan error)

	// Start a goroutine to generate keys asynchronously
	go func() {
		errCh <- r.generateKey()
		close(errCh) // Close the channel when done
	}()

	// Check for errors
	if err := <-errCh; err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// ReallocateKey moves a key back to the unused_keys set.
func (r *KeyGeneratorRepository) ReallocateKey(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Transactional move: Remove from used_keys, add to unused_keys
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, "used_keys", key)
		pipe.SAdd(ctx, "unused_keys", key)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to reallocate key: %w", err)
	}
	return nil
}

// generateKey creates and stores a unique key in unused_keys.
func (r *KeyGeneratorRepository) generateKey() error {
	const maxRetries = 10
	retries := 0

	for retries < maxRetries {
		key := make([]byte, 12) // Longer for uniqueness
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("key generation failed: %w", err)
		}

		encodedKey := base64.URLEncoding.EncodeToString(key)[:8]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Ensure key uniqueness
		if isUnused, err := r.isMemberOfSet(encodedKey, "unused_keys"); err != nil {
			return err
		} else if !isUnused {
			if isUsed, err := r.isMemberOfSet(encodedKey, "used_keys"); err != nil {
				return err
			} else if !isUsed {
				// Add key if unique
				if err := r.client.SAdd(ctx, "unused_keys", encodedKey).Err(); err != nil {
					return fmt.Errorf("failed to store new key: %w", err)
				}
				return nil
			}
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
