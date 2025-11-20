package factories

import "golang-social-media/apps/ecommerce-service/internal/domain/order"

// OrderFactory defines the contract for creating Order aggregates
type OrderFactory interface {
	CreateOrder(userID string) (*order.Order, error)
}


