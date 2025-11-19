package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// GetProductQuery retrieves a product by ID
type GetProductQuery interface {
	Execute(ctx context.Context, productID string) (product.Product, error)
}

