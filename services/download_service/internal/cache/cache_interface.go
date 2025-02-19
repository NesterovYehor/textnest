package cache

import (
	"context"

	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
)

// Cache defines the behavior for a cache implementation.
type Cache interface {
	// Set adds a key-value pair to the cache with an optional expiration time.
	Set(ctx context.Context, key string, metadata *pb.Metadata) error

	// Get retrieves a value from the cache by its key.
	// Returns the value and a boolean indicating if the key was found.
	Get(ctx context.Context, key string) (*pb.Metadata, bool, error)

	// Delete removes a key-value pair from the cache by its key.
	Delete(ctx context.Context, key string) error

	// Clear removes all entries from the cache.
	Clear(ctx context.Context) error

	// Close cache conection
	Close() error
}
