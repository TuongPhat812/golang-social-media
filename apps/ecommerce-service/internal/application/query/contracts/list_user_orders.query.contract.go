package contracts

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// ListUserOrdersQuery retrieves orders for a user
type ListUserOrdersQuery interface {
	Execute(ctx context.Context, userID string, limit int) ([]order.Order, error)
}

