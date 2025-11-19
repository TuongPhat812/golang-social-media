package order

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
	Version() int
	AggregateID() string
	AggregateType() string
}

// OrderCreatedEvent is a domain event emitted when an order is created
type OrderCreatedEvent struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	CreatedAt   string
	version     int // Event version
}

func (e OrderCreatedEvent) Type() string {
	return "OrderCreated"
}

func (e OrderCreatedEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e OrderCreatedEvent) AggregateID() string {
	return e.OrderID
}

func (e OrderCreatedEvent) AggregateType() string {
	return "Order"
}

// OrderItemAddedEvent is a domain event emitted when an item is added to an order
type OrderItemAddedEvent struct {
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
	SubTotal  float64
	UpdatedAt string
	version   int // Event version
}

func (e OrderItemAddedEvent) Type() string {
	return "OrderItemAdded"
}

func (e OrderItemAddedEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e OrderItemAddedEvent) AggregateID() string {
	return e.OrderID
}

func (e OrderItemAddedEvent) AggregateType() string {
	return "Order"
}

// OrderConfirmedEvent is a domain event emitted when an order is confirmed
type OrderConfirmedEvent struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	ItemCount   int
	ConfirmedAt string
	version     int // Event version
}

func (e OrderConfirmedEvent) Type() string {
	return "OrderConfirmed"
}

func (e OrderConfirmedEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e OrderConfirmedEvent) AggregateID() string {
	return e.OrderID
}

func (e OrderConfirmedEvent) AggregateType() string {
	return "Order"
}

// OrderCancelledEvent is a domain event emitted when an order is cancelled
type OrderCancelledEvent struct {
	OrderID    string
	UserID     string
	CancelledAt string
	version    int // Event version
}

func (e OrderCancelledEvent) Type() string {
	return "OrderCancelled"
}

func (e OrderCancelledEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e OrderCancelledEvent) AggregateID() string {
	return e.OrderID
}

func (e OrderCancelledEvent) AggregateType() string {
	return "Order"
}

