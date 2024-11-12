package keymanager_test

import (
	"fmt"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/keymanager"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var client *redis.Client

// Initialize the Redis client before running tests
func TestMain(m *testing.M) {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Ensure Redis is running on this address
	})

	// Ensure Redis connection is available
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	// Run tests
	m.Run()
}

// Test GetKey function
func TestGetKey(t *testing.T) {
	// Set up the initial state for testing
	client.SAdd(context.Background(), "unused_keys", "key123") // Add a test key to "unused_keys"

	// Call the function
	key, err := keymanager.GetKey(client)
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}

	// Check if the key was moved to "used_keys"
	isUsed, err := client.SIsMember(context.Background(), "used_keys", key).Result()
	if err != nil {
		t.Fatalf("Failed to check used_keys set: %v", err)
	}
	if !isUsed {
		t.Errorf("Expected key to be in 'used_keys' but it wasn't")
	}
}

// Test ReallocateKey function
func TestReallocateKey(t *testing.T) {
	// Set up initial keys
	client.SAdd(context.Background(), "used_keys", "key123") // Add a test key to "used_keys"

	// Call the function
	err := keymanager.ReallocateKey("key123", client)
	if err != nil {
		t.Fatalf("ReallocateKey failed: %v", err)
	}

	// Check if the key was moved back to "unused_keys"
	isUnused, err := client.SIsMember(context.Background(), "unused_keys", "key123").Result()
	if err != nil {
		t.Fatalf("Failed to check unused_keys set: %v", err)
	}
	if !isUnused {
		t.Errorf("Expected key to be in 'unused_keys' but it wasn't")
	}
}

// Test generateKey function
func TestGenerateKey(t *testing.T) {
	// Clean up before test
	client.Del(context.Background(), "unused_keys", "used_keys")

	// Call the function to generate a new key
	err := keymanager.GenerateKey(client)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	// Verify the key was added to "unused_keys"
	keys, err := client.SMembers(context.Background(), "unused_keys").Result()
	if err != nil {
		t.Fatalf("Failed to get members of unused_keys: %v", err)
	}
	if len(keys) == 0 {
		t.Errorf("Expected to have at least one key in 'unused_keys', but found none")
	}
}

// Test FillKeys function
func TestFillKeys(t *testing.T) {
	// Clean up before test
	client.Del(context.Background(), "unused_keys")

	// Ensure there are less than the threshold keys
	client.SAdd(context.Background(), "unused_keys", "key1", "key2")

	// Call the function
	keymanager.FillKeys(client, 5) // Threshold set to 5

	// Verify that there are now 5 keys in "unused_keys"
	count, err := client.SCard(context.Background(), "unused_keys").Result()
	if err != nil {
		t.Fatalf("Failed to get count of unused_keys: %v", err)
	}
	if count < 5 {
		t.Errorf("Expected at least 5 keys in 'unused_keys', but found %d", count)
	}
}

// Test IsKeyValid function
func TestIsKeyValid(t *testing.T) {
	v := validator.NewValidator()

	// Test valid key
	keymanager.IsKeyValid(v, "validkey")
	if len(v.Errors) > 0 {
		t.Errorf("Expected no validation errors, but got: %v", v.Errors)
	}

	// Test invalid key (empty)
	keymanager.IsKeyValid(v, "")
	if len(v.Errors) == 0 {
		t.Errorf("Expected validation errors, but got none")
	}

	// Test invalid key (incorrect length)
	keymanager.IsKeyValid(v, "short")
	if len(v.Errors) == 0 {
		t.Errorf("Expected validation errors, but got none")
	}
}
