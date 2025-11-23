package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	"golang-social-media/pkg/cache"
)

// UserCache handles caching for User entities in chat service
type UserCache struct {
	cache cache.Cache
	ttl   time.Duration
}

// NewUserCache creates a new UserCache
func NewUserCache(cache cache.Cache) *UserCache {
	return &UserCache{
		cache: cache,
		ttl:   30 * time.Minute, // User cache TTL
	}
}

// GetUser retrieves a user from cache
func (c *UserCache) GetUser(ctx context.Context, id string) (*persistence.UserModel, error) {
	key := c.userKey(id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var user persistence.UserModel
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUser stores a user in cache
func (c *UserCache) SetUser(ctx context.Context, user *persistence.UserModel) error {
	key := c.userKey(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

// DeleteUser removes a user from cache
func (c *UserCache) DeleteUser(ctx context.Context, id string) error {
	key := c.userKey(id)
	return c.cache.Delete(ctx, key)
}

// userKey generates a cache key for a user
func (c *UserCache) userKey(id string) string {
	return fmt.Sprintf("chat:user:%s", id)
}

