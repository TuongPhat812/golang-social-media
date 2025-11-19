package product

import (
	"errors"
	"strings"
	"time"
)

// Status represents product status
type Status string

const (
	StatusActive     Status = "active"
	StatusInactive   Status = "inactive"
	StatusOutOfStock Status = "out_of_stock"
)

// Product represents a product entity
// This is a standalone entity (not an aggregate root)
// It has its own identity and business logic
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64 // In production, should use Money value object
	Stock       int
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the product
func (p Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}

// Create is a domain method that creates a product and adds a domain event
func (p *Product) Create() {
	p.addEvent(ProductCreatedEvent{
		ProductID:   p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	})
}

// UpdateStock updates the product stock and adds a domain event
func (p *Product) UpdateStock(newStock int) error {
	if newStock < 0 {
		return errors.New("stock cannot be negative")
	}

	oldStock := p.Stock
	p.Stock = newStock
	p.UpdatedAt = time.Now().UTC()

	// Update status based on stock
	if p.Stock == 0 {
		p.Status = StatusOutOfStock
	} else if p.Status == StatusOutOfStock {
		p.Status = StatusActive
	}

	// Add domain event
	p.addEvent(ProductStockUpdatedEvent{
		ProductID: p.ID,
		OldStock:  oldStock,
		NewStock:  newStock,
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	})

	return nil
}

// DecreaseStock decreases stock by quantity (used when order is confirmed)
func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if p.Stock < quantity {
		return errors.New("insufficient stock")
	}

	return p.UpdateStock(p.Stock - quantity)
}

// IncreaseStock increases stock by quantity (used when order is cancelled)
func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	return p.UpdateStock(p.Stock + quantity)
}

// IsAvailable returns true if product is available for purchase
func (p Product) IsAvailable() bool {
	return p.Status == StatusActive && p.Stock > 0
}

// Events returns all domain events
func (p Product) Events() []DomainEvent {
	return p.events
}

// ClearEvents clears all domain events
func (p *Product) ClearEvents() {
	p.events = nil
}

// addEvent adds a domain event (internal method)
func (p *Product) addEvent(event DomainEvent) {
	p.events = append(p.events, event)
}

