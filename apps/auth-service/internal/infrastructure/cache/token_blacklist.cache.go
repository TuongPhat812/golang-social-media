package cache

import (
	"context"
	"fmt"
	"time"

	"golang-social-media/pkg/cache"
)

// TokenBlacklist handles token blacklist for logout functionality
type TokenBlacklist struct {
	cache cache.Cache
	ttl   time.Duration
}

// NewTokenBlacklist creates a new TokenBlacklist
func NewTokenBlacklist(cache cache.Cache) *TokenBlacklist {
	return &TokenBlacklist{
		cache: cache,
		ttl:   24 * time.Hour, // Blacklist tokens for 24 hours (should match max token expiration)
	}
}

// BlacklistToken adds a token to the blacklist
func (b *TokenBlacklist) BlacklistToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := b.tokenKey(tokenID)
	// Use token expiration or default TTL, whichever is longer
	ttl := expiration
	if ttl < b.ttl {
		ttl = b.ttl
	}
	return b.cache.Set(ctx, key, []byte("blacklisted"), ttl)
}

// IsBlacklisted checks if a token is blacklisted
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := b.tokenKey(tokenID)
	exists, err := b.cache.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// tokenKey generates a cache key for a token
func (b *TokenBlacklist) tokenKey(tokenID string) string {
	return fmt.Sprintf("auth:token:blacklist:%s", tokenID)
}

