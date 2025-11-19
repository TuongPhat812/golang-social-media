package contracts

import (
	"context"
)

// Subscriber represents a generic subscriber interface
type Subscriber interface {
	Consume(ctx context.Context)
	Close() error
}

// ProductCreatedSubscriber subscribes to ProductCreated events
type ProductCreatedSubscriber interface {
	Subscriber
}

// ProductStockUpdatedSubscriber subscribes to ProductStockUpdated events
type ProductStockUpdatedSubscriber interface {
	Subscriber
}

// OrderCreatedSubscriber subscribes to OrderCreated events
type OrderCreatedSubscriber interface {
	Subscriber
}

// OrderItemAddedSubscriber subscribes to OrderItemAdded events
type OrderItemAddedSubscriber interface {
	Subscriber
}

// OrderConfirmedSubscriber subscribes to OrderConfirmed events
type OrderConfirmedSubscriber interface {
	Subscriber
}

// OrderCancelledSubscriber subscribes to OrderCancelled events
type OrderCancelledSubscriber interface {
	Subscriber
}

