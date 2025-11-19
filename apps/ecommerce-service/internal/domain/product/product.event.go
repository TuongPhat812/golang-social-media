package product

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
	Version() int
	AggregateID() string
	AggregateType() string
}

// ProductCreatedEvent is a domain event emitted when a product is created
type ProductCreatedEvent struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   string
	version     int // Event version
}

func (e ProductCreatedEvent) Type() string {
	return "ProductCreated"
}

func (e ProductCreatedEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e ProductCreatedEvent) AggregateID() string {
	return e.ProductID
}

func (e ProductCreatedEvent) AggregateType() string {
	return "Product"
}

// ProductStockUpdatedEvent is a domain event emitted when product stock is updated
type ProductStockUpdatedEvent struct {
	ProductID string
	OldStock  int
	NewStock  int
	UpdatedAt string
	version   int // Event version
}

func (e ProductStockUpdatedEvent) Type() string {
	return "ProductStockUpdated"
}

func (e ProductStockUpdatedEvent) Version() int {
	if e.version == 0 {
		return 1 // Default version
	}
	return e.version
}

func (e ProductStockUpdatedEvent) AggregateID() string {
	return e.ProductID
}

func (e ProductStockUpdatedEvent) AggregateType() string {
	return "Product"
}

