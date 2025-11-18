package contracts

import (
	"context"
)

// Publisher defines the interface for publishing events
type Publisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
	Close() error
}

