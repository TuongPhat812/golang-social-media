package contracts

import (
	"context"
)

// Subscriber represents a generic subscriber interface
type Subscriber interface {
	Consume(ctx context.Context)
	Close() error
}
