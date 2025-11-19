package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// GetOrderQuery retrieves an order by ID
type GetOrderQuery interface {
	Execute(ctx context.Context, orderID string) (order.Order, error)
}

