package products

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// Repository defines the interface for product persistence
type Repository interface {
	Create(ctx context.Context, p *product.Product) error
	FindByID(ctx context.Context, id string) (product.Product, error)
	Update(ctx context.Context, p *product.Product) error
	List(ctx context.Context, status *product.Status, limit, offset int) ([]product.Product, error)
}

