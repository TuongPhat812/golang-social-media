package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// OrderCache handles caching for Order entities
type OrderCache struct {
	cache   Cache
	ttl     time.Duration
	listTTL time.Duration
}

// NewOrderCache creates a new OrderCache
func NewOrderCache(cache Cache) *OrderCache {
	return &OrderCache{
		cache:   cache,
		ttl:     10 * time.Minute,      // Individual order cache TTL
		listTTL: 3 * time.Minute,       // Order list cache TTL
	}
}

// GetOrder retrieves an order from cache
func (c *OrderCache) GetOrder(ctx context.Context, id string) (*order.Order, error) {
	key := c.orderKey(id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var o order.Order
	if err := json.Unmarshal(data, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// SetOrder stores an order in cache
func (c *OrderCache) SetOrder(ctx context.Context, o *order.Order) error {
	key := c.orderKey(o.ID)
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

// DeleteOrder removes an order from cache
func (c *OrderCache) DeleteOrder(ctx context.Context, id string) error {
	key := c.orderKey(id)
	return c.cache.Delete(ctx, key)
}

// GetOrderList retrieves an order list from cache
func (c *OrderCache) GetOrderList(ctx context.Context, userID string, limit int) ([]order.Order, error) {
	key := c.orderListKey(userID, limit)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var orders []order.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// SetOrderList stores an order list in cache
func (c *OrderCache) SetOrderList(ctx context.Context, userID string, limit int, orders []order.Order) error {
	key := c.orderListKey(userID, limit)
	data, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.listTTL)
}

// InvalidateOrderList invalidates all order list caches for a user
func (c *OrderCache) InvalidateOrderList(ctx context.Context, userID string) error {
	pattern := c.orderListPattern(userID)
	return c.cache.DeletePattern(ctx, pattern)
}

// orderKey generates a cache key for an order
func (c *OrderCache) orderKey(id string) string {
	return fmt.Sprintf("order:%s", id)
}

// orderListKey generates a cache key for an order list
func (c *OrderCache) orderListKey(userID string, limit int) string {
	return fmt.Sprintf("order:list:%s:%d", userID, limit)
}

// orderListPattern generates a pattern for order list keys
func (c *OrderCache) orderListPattern(userID string) string {
	return fmt.Sprintf("order:list:%s:*", userID)
}
