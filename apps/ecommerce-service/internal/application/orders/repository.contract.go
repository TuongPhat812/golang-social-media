package orders

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
)

// Repository defines the interface for order persistence
type Repository interface {
	Create(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, id string) (order.Order, error)
	Update(ctx context.Context, o *order.Order) error
	ListByUser(ctx context.Context, userID string, limit int) ([]order.Order, error)
}

