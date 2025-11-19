package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// ProductCache handles caching for Product entities
type ProductCache struct {
	cache   Cache
	ttl     time.Duration
	listTTL time.Duration
}

// NewProductCache creates a new ProductCache
func NewProductCache(cache Cache) *ProductCache {
	return &ProductCache{
		cache:   cache,
		ttl:     15 * time.Minute,      // Individual product cache TTL
		listTTL: 5 * time.Minute,       // Product list cache TTL
	}
}

// GetProduct retrieves a product from cache
func (c *ProductCache) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	key := c.productKey(id)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var p product.Product
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// SetProduct stores a product in cache
func (c *ProductCache) SetProduct(ctx context.Context, p *product.Product) error {
	key := c.productKey(p.ID)
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.ttl)
}

// DeleteProduct removes a product from cache
func (c *ProductCache) DeleteProduct(ctx context.Context, id string) error {
	key := c.productKey(id)
	return c.cache.Delete(ctx, key)
}

// GetProductList retrieves a product list from cache
func (c *ProductCache) GetProductList(ctx context.Context, status *product.Status, limit, offset int) ([]product.Product, error) {
	key := c.productListKey(status, limit, offset)
	data, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var products []product.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// SetProductList stores a product list in cache
func (c *ProductCache) SetProductList(ctx context.Context, status *product.Status, limit, offset int, products []product.Product) error {
	key := c.productListKey(status, limit, offset)
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, data, c.listTTL)
}

// InvalidateProductList invalidates all product list caches
func (c *ProductCache) InvalidateProductList(ctx context.Context) error {
	pattern := "product:list:*"
	return c.cache.DeletePattern(ctx, pattern)
}

// productKey generates a cache key for a product
func (c *ProductCache) productKey(id string) string {
	return fmt.Sprintf("product:%s", id)
}

// productListKey generates a cache key for a product list
func (c *ProductCache) productListKey(status *product.Status, limit, offset int) string {
	statusStr := "all"
	if status != nil {
		statusStr = string(*status)
	}
	return fmt.Sprintf("product:list:%s:%d:%d", statusStr, limit, offset)
}
