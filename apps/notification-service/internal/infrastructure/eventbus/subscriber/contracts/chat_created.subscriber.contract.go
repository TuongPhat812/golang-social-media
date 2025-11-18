package contracts

import (
	"context"
)

// ChatCreatedSubscriber consumes ChatCreated events
type ChatCreatedSubscriber interface {
	Consume(ctx context.Context)
	Close() error
}

