package contracts

import (
	"context"
)

// CancelOrderCommand cancels an order
type CancelOrderCommand interface {
	Execute(ctx context.Context, orderID string) error
}

