package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	log    *zerolog.Logger
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Component("ecommerce.cache").
		Info().
		Str("addr", addr).
		Int("db", db).
		Msg("redis cache connected")

	return &RedisCache{
		client: client,
		log:    logger.Component("ecommerce.cache"),
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		c.log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to get from cache")
		return nil, err
	}
	return val, nil
}

// Set stores a value in cache with expiration
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := c.client.Set(ctx, key, value, expiration).Err(); err != nil {
		c.log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to set cache")
		return err
	}
	return nil
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to delete from cache")
		return err
	}
	return nil
}

// DeletePattern removes all keys matching a pattern
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		c.log.Error().
			Err(err).
			Str("pattern", pattern).
			Msg("failed to scan keys")
		return err
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			c.log.Error().
				Err(err).
				Str("pattern", pattern).
				Int("key_count", len(keys)).
				Msg("failed to delete keys by pattern")
			return err
		}
	}

	return nil
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to check cache existence")
		return false, err
	}
	return count > 0, nil
}

// Close closes the cache connection
func (c *RedisCache) Close() error {
	if err := c.client.Close(); err != nil {
		c.log.Error().
			Err(err).
			Msg("failed to close redis cache")
		return err
	}
	return nil
}

// ErrCacheMiss is returned when a key is not found in cache
var ErrCacheMiss = &CacheError{Message: "cache miss"}

// CacheError represents a cache error
type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

