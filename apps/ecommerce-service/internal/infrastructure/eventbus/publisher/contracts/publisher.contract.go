package contracts

import (
	"context"
)

// EcommercePublisher defines the contract for publishing ecommerce events
type EcommercePublisher interface {
	PublishProductCreated(ctx context.Context, event ProductCreated) error
	PublishProductStockUpdated(ctx context.Context, event ProductStockUpdated) error
	PublishOrderCreated(ctx context.Context, event OrderCreated) error
	PublishOrderItemAdded(ctx context.Context, event OrderItemAdded) error
	PublishOrderConfirmed(ctx context.Context, event OrderConfirmed) error
	PublishOrderCancelled(ctx context.Context, event OrderCancelled) error
	Close() error
}

// ProductCreated represents a product created event
type ProductCreated struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   string
}

// ProductStockUpdated represents a product stock updated event
type ProductStockUpdated struct {
	ProductID string
	OldStock  int
	NewStock  int
	UpdatedAt string
}

// OrderCreated represents an order created event
type OrderCreated struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	CreatedAt   string
}

// OrderItemAdded represents an order item added event
type OrderItemAdded struct {
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
	SubTotal  float64
	UpdatedAt string
}

// OrderConfirmed represents an order confirmed event
type OrderConfirmed struct {
	OrderID     string
	UserID      string
	TotalAmount  float64
	ItemCount   int
	ConfirmedAt string
}

// OrderCancelled represents an order cancelled event
type OrderCancelled struct {
	OrderID    string
	UserID     string
	CancelledAt string
}

