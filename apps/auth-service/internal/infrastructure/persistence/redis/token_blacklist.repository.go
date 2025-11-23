package redis

import (
	"context"
	"fmt"
	"time"

	"golang-social-media/pkg/cache"
	"golang-social-media/pkg/logger"
)

// TokenBlacklistRepository manages blacklisted tokens in Redis
type TokenBlacklistRepository struct {
	cache cache.Cache
}

// NewTokenBlacklistRepository creates a new token blacklist repository
func NewTokenBlacklistRepository(cache cache.Cache) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{
		cache: cache,
	}
}

// AddToken adds a token to the blacklist until expiration
func (r *TokenBlacklistRepository) AddToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := r.tokenKey(tokenID)
	// Store token ID with expiration (TTL)
	return r.cache.Set(ctx, key, []byte("blacklisted"), expiration)
}

// IsBlacklisted checks if a token is blacklisted
func (r *TokenBlacklistRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := r.tokenKey(tokenID)
	exists, err := r.cache.Exists(ctx, key)
	if err != nil {
		logger.Component("auth.repository.token_blacklist").
			Error().
			Err(err).
			Str("token_id", tokenID).
			Msg("failed to check token blacklist")
		return false, err
	}
	return exists, nil
}

// RemoveToken removes a token from blacklist (if needed)
func (r *TokenBlacklistRepository) RemoveToken(ctx context.Context, tokenID string) error {
	key := r.tokenKey(tokenID)
	return r.cache.Delete(ctx, key)
}

// tokenKey generates a cache key for a token
func (r *TokenBlacklistRepository) tokenKey(tokenID string) string {
	return fmt.Sprintf("auth:token:blacklist:%s", tokenID)
}

