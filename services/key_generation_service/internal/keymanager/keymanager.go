package keymanager

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/redis/go-redis/v9"
)

const redisTimeout = time.Second * 5 // Set a consistent timeout duration

func StartRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func GetKey(rdb *redis.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	// Retrieve a random key from Redis
	key, err := rdb.SRandMember(ctx, "unused_keys").Result()
	if err != nil {
		return "", err
	}

	// Step 2: Remove the key from the unused_keys set
	_, err = rdb.SRem(ctx, "unused_keys", key).Result()
	if err != nil {
		return "", err
	}

	// Step 3: Add the key to the used_keys set
	_, err = rdb.SAdd(ctx, "used_keys", key).Result()
	if err != nil {
		return "", err
	}

	// Start a goroutine to generate a new key and store it in Redis
	go func() {
		if err := generateKey(rdb); err != nil {
			log.Println("Error generating key:", err)
		}
	}()

	return key, nil
}

func ReallocateKey(key string, rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	// Remove the key from the used_keys set
	_, err := rdb.SRem(ctx, "used_keys", key).Result()
	if err != nil {
		return err // Return error if removal fails
	}

	// Add the key to the unused_keys set
	_, err = rdb.SAdd(ctx, "unused_keys", key).Result()
	if err != nil {
		return err // Return error if addition fails
	}

	return nil // Return nil if everything is successful
}

func generateKey(rdb *redis.Client) error {
	key := make([]byte, 12) // Generate a longer key for more uniqueness
	for {
		if _, err := rand.Read(key); err != nil {
			return err
		}

		// Encode and slice the key to the desired length
		encodedKey := base64.URLEncoding.EncodeToString(key)[:8]
		ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
		defer cancel()

		// Check if the key exists in either set
		if isUnused, err := isMemberOfSet(rdb, encodedKey, "unused_keys"); err != nil {
			return err
		} else if !isUnused {
			if isUsed, err := isMemberOfSet(rdb, encodedKey, "used_keys"); err != nil {
				return err
			} else if !isUsed {
				// The key is unique, store it in Redis
				err := rdb.SAdd(ctx, "unused_keys", encodedKey).Err()
				if err != nil {
					return errors.New("Failed to store new key: " + encodedKey)
				}
				return nil // Successfully generated and stored unique key
			}
		}
	}
}

func isMemberOfSet(rdb *redis.Client, value string, setName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	// Check if the value is a member of the specified set
	isMember, err := rdb.SIsMember(ctx, setName, value).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func IsKeyValid(v *validator.Validator, key string) {
	v.Check(key == "", "key", "key must be provided")
	v.Check(len([]rune(key)) != 8, "key", "key must be 8 chars long")
}
