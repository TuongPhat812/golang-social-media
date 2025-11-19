package product

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// ProductCreatedEvent is a domain event emitted when a product is created
type ProductCreatedEvent struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   string
}

func (e ProductCreatedEvent) Type() string {
	return "ProductCreated"
}

// ProductStockUpdatedEvent is a domain event emitted when product stock is updated
type ProductStockUpdatedEvent struct {
	ProductID string
	OldStock  int
	NewStock  int
	UpdatedAt string
}

func (e ProductStockUpdatedEvent) Type() string {
	return "ProductStockUpdated"
}

