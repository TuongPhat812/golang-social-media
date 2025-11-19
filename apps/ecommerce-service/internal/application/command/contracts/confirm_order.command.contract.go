package contracts

import (
	"context"
)

// ConfirmOrderCommand confirms an order
type ConfirmOrderCommand interface {
	Execute(ctx context.Context, orderID string) error
}

