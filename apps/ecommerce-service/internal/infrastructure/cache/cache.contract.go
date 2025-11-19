package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with expiration
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// DeletePattern removes all keys matching a pattern
	DeletePattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// Close closes the cache connection
	Close() error
}

