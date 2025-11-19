package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// AddOrderItemCommand adds an item to an order
type AddOrderItemCommand interface {
	Execute(ctx context.Context, req AddOrderItemCommandRequest) (order.Order, error)
}

// AddOrderItemCommandRequest represents the request for adding an item to an order
type AddOrderItemCommandRequest struct {
	OrderID   string
	ProductID string
	Quantity  int
}

