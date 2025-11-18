package contracts

import (
	"context"
)

// Query represents a read operation that does not modify state
type Query interface {
	Execute(ctx context.Context) (interface{}, error)
}

