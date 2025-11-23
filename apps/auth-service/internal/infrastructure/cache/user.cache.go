package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/pkg/cache"
)

// UserCache handles caching for User entities in auth service
type UserCache struct {
	cache cache.Cache
	ttl   time.Duration
}

// NewUserCache creates a new UserCache
func NewUserCache(cache cache.Cache) *UserCache {
	return &UserCache{
		cache: cache,
		ttl:   15 * time.Minute, // User cache TTL
	}
}

// GetUserByID retrieves a user from cache by ID
func (c *UserCache) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	key := c.userKey(id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var u user.User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail retrieves a user by email from cache
func (c *UserCache) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	key := c.userEmailKey(email)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var u user.User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// SetUser stores a user in cache (by both ID and email)
func (c *UserCache) SetUser(ctx context.Context, u *user.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	// Cache by ID
	keyID := c.userKey(u.ID)
	if err := c.cache.Set(ctx, keyID, data, c.ttl); err != nil {
		return err
	}

	// Cache by email
	keyEmail := c.userEmailKey(u.Email)
	return c.cache.Set(ctx, keyEmail, data, c.ttl)
}

// DeleteUser removes a user from cache by ID and email
func (c *UserCache) DeleteUser(ctx context.Context, id, email string) error {
	if err := c.cache.Delete(ctx, c.userKey(id)); err != nil {
		return err
	}
	return c.cache.Delete(ctx, c.userEmailKey(email))
}

// userKey generates a cache key for a user by ID
func (c *UserCache) userKey(id string) string {
	return fmt.Sprintf("auth:user:id:%s", id)
}

// userEmailKey generates a cache key for a user by email
func (c *UserCache) userEmailKey(email string) string {
	return fmt.Sprintf("auth:user:email:%s", email)
}
