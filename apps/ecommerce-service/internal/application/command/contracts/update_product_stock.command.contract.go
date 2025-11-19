package contracts

import (
	"context"
)

// UpdateProductStockCommand updates product stock
type UpdateProductStockCommand interface {
	Execute(ctx context.Context, productID string, newStock int) error
}

