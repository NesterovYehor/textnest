package testutils

import (
	"context"
	"net/url"
	"testing"

	"github.com/NesterovYehor/TextNest/pkg/test/container"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func StartTestRedis(t *testing.T, ctx context.Context) (*redis.Client, func()) {
	redisOpts := container.RedisContainerOpts{
		Addr: 9090,
	}

	// Start Redis container
	redisContainer, err := container.StartRedis(ctx, &redisOpts)
	require.NoError(t, err, "Failed to start Redis container")

	// Get Redis connection URI
	redisUri, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get Redis connection string")

	// Parse URI
	parsedUri, err := url.Parse(redisUri)
	require.NoError(t, err, "Failed to parse Redis URI")
	require.NotEmpty(t, parsedUri.Host, "Parsed URI host is empty")

	// Initialize Redis client
	client := redis.NewClient(&redis.Options{Addr: parsedUri.Host})
	require.NotNil(t, client, "Redis client initialization failed")

	// Cleanup function to terminate Redis container
	cleanup := func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate Redis container: %v", err)
		}
	}

	return client, cleanup
}
