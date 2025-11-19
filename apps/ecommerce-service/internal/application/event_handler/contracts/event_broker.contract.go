package contracts

import (
	"context"
)

// EventBrokerPublisher represents the contract for publishing events to the event broker
type EventBrokerPublisher interface {
	PublishProductCreated(ctx context.Context, payload ProductCreatedPayload) error
	PublishProductStockUpdated(ctx context.Context, payload ProductStockUpdatedPayload) error
	PublishOrderCreated(ctx context.Context, payload OrderCreatedPayload) error
	PublishOrderItemAdded(ctx context.Context, payload OrderItemAddedPayload) error
	PublishOrderConfirmed(ctx context.Context, payload OrderConfirmedPayload) error
	PublishOrderCancelled(ctx context.Context, payload OrderCancelledPayload) error
}

// ProductCreatedPayload represents the payload for ProductCreated event
type ProductCreatedPayload struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   string
}

// ProductStockUpdatedPayload represents the payload for ProductStockUpdated event
type ProductStockUpdatedPayload struct {
	ProductID string
	OldStock  int
	NewStock  int
	UpdatedAt string
}

// OrderCreatedPayload represents the payload for OrderCreated event
type OrderCreatedPayload struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	CreatedAt   string
}

// OrderItemAddedPayload represents the payload for OrderItemAdded event
type OrderItemAddedPayload struct {
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
	SubTotal  float64
	UpdatedAt string
}

// OrderConfirmedPayload represents the payload for OrderConfirmed event
type OrderConfirmedPayload struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	ConfirmedAt string
}

// OrderCancelledPayload represents the payload for OrderCancelled event
type OrderCancelledPayload struct {
	OrderID    string
	UserID     string
	CancelledAt string
}

