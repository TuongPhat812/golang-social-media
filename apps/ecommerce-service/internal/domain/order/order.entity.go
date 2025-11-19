package order

import (
	"errors"
	"time"
)

// Status represents order status
type Status string

const (
	StatusDraft     Status = "draft"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
	StatusCompleted Status = "completed"
)

// Order represents an order aggregate root
// Aggregate Root: Entry point to access the Order aggregate, manages OrderItems
type Order struct {
	ID          string
	UserID      string
	Status      Status
	Items       []OrderItem // Value Objects within aggregate
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the order
func (o Order) Validate() error {
	if o.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if len(o.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	if o.TotalAmount < 0 {
		return errors.New("total amount cannot be negative")
	}
	return nil
}

// Create is a domain method that creates an order and adds a domain event
func (o *Order) Create() {
	o.addEvent(OrderCreatedEvent{
		OrderID:     o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		ItemCount:   len(o.Items),
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
	})
}

// AddItem adds an item to the order
// This method enforces aggregate boundary - items can only be added through the aggregate root
func (o *Order) AddItem(item OrderItem) error {
	if o.Status != StatusDraft {
		return errors.New("can only add items to draft orders")
	}

	// Validate the value object
	if err := item.Validate(); err != nil {
		return err
	}

	o.Items = append(o.Items, item)
	o.recalculateTotal()
	o.UpdatedAt = time.Now().UTC()

	// Add domain event
	o.addEvent(OrderItemAddedEvent{
		OrderID:   o.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		SubTotal:  item.SubTotal,
		UpdatedAt: o.UpdatedAt.Format(time.RFC3339),
	})

	return nil
}

// Confirm confirms the order and adds a domain event
func (o *Order) Confirm() error {
	if o.Status != StatusDraft {
		return errors.New("can only confirm draft orders")
	}
	if len(o.Items) == 0 {
		return errors.New("cannot confirm order with no items")
	}

	o.Status = StatusConfirmed
	o.UpdatedAt = time.Now().UTC()

	// Add domain event
	o.addEvent(OrderConfirmedEvent{
		OrderID:     o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		ItemCount:   len(o.Items),
		ConfirmedAt: o.UpdatedAt.Format(time.RFC3339),
	})

	return nil
}

// Cancel cancels the order and adds a domain event
func (o *Order) Cancel() error {
	if o.Status == StatusCancelled {
		return errors.New("order is already cancelled")
	}
	if o.Status == StatusCompleted {
		return errors.New("cannot cancel completed orders")
	}

	o.Status = StatusCancelled
	o.UpdatedAt = time.Now().UTC()

	// Add domain event
	o.addEvent(OrderCancelledEvent{
		OrderID:    o.ID,
		UserID:     o.UserID,
		CancelledAt: o.UpdatedAt.Format(time.RFC3339),
	})

	return nil
}

// recalculateTotal recalculates the total amount based on items
func (o *Order) recalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.SubTotal
	}
	o.TotalAmount = total
}

// Events returns all domain events
func (o Order) Events() []DomainEvent {
	return o.events
}

// ClearEvents clears all domain events
func (o *Order) ClearEvents() {
	o.events = nil
}

// addEvent adds a domain event (internal method)
func (o *Order) addEvent(event DomainEvent) {
	o.events = append(o.events, event)
}

