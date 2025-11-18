package contracts

import (
	"context"
)

// Subscriber defines the interface for consuming events
type Subscriber interface {
	Consume(ctx context.Context)
	Close() error
}

