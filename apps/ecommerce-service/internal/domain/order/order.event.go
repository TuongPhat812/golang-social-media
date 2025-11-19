package order

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// OrderCreatedEvent is a domain event emitted when an order is created
type OrderCreatedEvent struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	CreatedAt   string
}

func (e OrderCreatedEvent) Type() string {
	return "OrderCreated"
}

// OrderItemAddedEvent is a domain event emitted when an item is added to an order
type OrderItemAddedEvent struct {
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
	SubTotal  float64
	UpdatedAt string
}

func (e OrderItemAddedEvent) Type() string {
	return "OrderItemAdded"
}

// OrderConfirmedEvent is a domain event emitted when an order is confirmed
type OrderConfirmedEvent struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	ConfirmedAt string
}

func (e OrderConfirmedEvent) Type() string {
	return "OrderConfirmed"
}

// OrderCancelledEvent is a domain event emitted when an order is cancelled
type OrderCancelledEvent struct {
	OrderID    string
	UserID     string
	CancelledAt string
}

func (e OrderCancelledEvent) Type() string {
	return "OrderCancelled"
}

