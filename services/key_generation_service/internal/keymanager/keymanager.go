package keymanager

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

const timeout = time.Second * 5 // Set a consistent timeout duration

func GetKey(client *redis.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Retrieve a random key from Redis
	key, err := client.SRandMember(ctx, "unused_keys").Result()
	if err != nil {
		return "", err
	}

	// Step 2: Remove the key from the unused_keys set
	_, err = client.SRem(ctx, "unused_keys", key).Result()
	if err != nil {
		return "", err
	}

	// Step 3: Add the key to the used_keys set
	_, err = client.SAdd(ctx, "used_keys", key).Result()
	if err != nil {
		return "", err
	}

	// Start a goroutine to generate a new key and store it in Redis
	go func() {
		if err := generateKey(client); err != nil {
			log.Println("Error generating key:", err)
		}
	}()

	return key, nil
}

func ReallocateKey(key string, client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Remove the key from the used_keys set
	_, err := client.SRem(ctx, "used_keys", key).Result()
	if err != nil {
		return err // Return error if removal fails
	}

	// Add the key to the unused_keys set
	_, err = client.SAdd(ctx, "unused_keys", key).Result()
	if err != nil {
		return err // Return error if addition fails
	}

	return nil // Return nil if everything is successful
}

func generateKey(client *redis.Client) error {
	const maxRetries = 10 // Retry limit
	retries := 0

	key := make([]byte, 12) // Generate a longer key for more uniqueness
	for retries < maxRetries {
		if _, err := rand.Read(key); err != nil {
			return err
		}

		// Encode and slice the key to the desired length
		encodedKey := base64.URLEncoding.EncodeToString(key)[:8]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Check if the key exists in either set
		if isUnused, err := isMemberOfSet(client, encodedKey, "unused_keys"); err != nil {
			return err
		} else if !isUnused {
			if isUsed, err := isMemberOfSet(client, encodedKey, "used_keys"); err != nil {
				return err
			} else if !isUsed {
				// The key is unique, store it in Redis
				err := client.SAdd(ctx, "unused_keys", encodedKey).Err()
				if err != nil {
					return fmt.Errorf("Failed to store new key: %s", encodedKey)
				}
				return nil // Successfully generated and stored unique key
			}
		}
		retries++
	}

	return errors.New("max retries reached for key generation")
}

func isMemberOfSet(client *redis.Client, value string, setName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Check if the value is a member of the specified set
	isMember, err := client.SIsMember(ctx, setName, value).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func FillKeys(client *redis.Client, threshold int64) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := client.SCard(ctx, "unused_keys").Result() // Count of unused keys
	if err != nil {
		log.Fatalf("Failed to get count of unused keys: %v", err)
	}

	for count < threshold {
		err := generateKey(client) // Ensure generateKey returns valid key and error
		if err != nil {
			log.Fatalf("Failed to generate a new key: %v", err)
		}

		count++
	}
}

func IsKeyValid(v *validator.Validator, key string) {
	v.Check(key == "", "key", "key must be provided")
	v.Check(len([]rune(key)) != 8, "key", "key must be 8 chars long")
}
