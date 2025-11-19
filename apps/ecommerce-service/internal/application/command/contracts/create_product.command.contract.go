package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// CreateProductCommand creates a new product
type CreateProductCommand interface {
	Execute(ctx context.Context, req CreateProductCommandRequest) (product.Product, error)
}

// CreateProductCommandRequest represents the request for creating a product
type CreateProductCommandRequest struct {
	Name        string
	Description string
	Price       float64
	Stock       int
}

