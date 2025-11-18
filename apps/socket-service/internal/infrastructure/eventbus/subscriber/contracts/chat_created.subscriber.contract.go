package contracts

import (
	"context"
)

// ChatCreatedSubscriber subscribes to ChatCreated events
type ChatCreatedSubscriber interface {
	Consume(ctx context.Context)
	Close() error
}
