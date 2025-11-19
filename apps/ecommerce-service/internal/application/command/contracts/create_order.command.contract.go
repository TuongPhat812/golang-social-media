package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// CreateOrderCommand creates a new order
type CreateOrderCommand interface {
	Execute(ctx context.Context, req CreateOrderCommandRequest) (order.Order, error)
}

// CreateOrderCommandRequest represents the request for creating an order
type CreateOrderCommandRequest struct {
	UserID string
}

