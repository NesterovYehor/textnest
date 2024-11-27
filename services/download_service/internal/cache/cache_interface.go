package cache

import (
	"time"

)

// Cache defines the behavior for a cache implementation.
type Cache interface {
	// Set adds a key-value pair to the cache with an optional expiration time.
	Set(key string, value []byte, expiration time.Duration) error

	// Get retrieves a value from the cache by its key.
	// Returns the value and a boolean indicating if the key was found.
	Get(key string) ([]byte, bool, error)

	// Delete removes a key-value pair from the cache by its key.
	Delete(key string) error

	// Clear removes all entries from the cache.
	Clear() error

	// Close cache conection
	Close() error
}
