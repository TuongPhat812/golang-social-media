package contracts

import (
	"context"
)

// Command represents a write operation that modifies state
type Command interface {
	Execute(ctx context.Context) error
}

