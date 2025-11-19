package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// ListProductsQuery retrieves a list of products
type ListProductsQuery interface {
	Execute(ctx context.Context, req ListProductsQueryRequest) ([]product.Product, error)
}

// ListProductsQueryRequest represents the request for listing products
type ListProductsQueryRequest struct {
	Status string // Optional filter by status
	Limit  int    // Pagination limit
	Offset int    // Pagination offset
}

