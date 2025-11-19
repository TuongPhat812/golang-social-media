package order

import (
	"errors"
)

// OrderItem represents an item in an order (Value Object)
// Value Objects are immutable and defined by their values, not identity
type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
	SubTotal  float64
}

// NewOrderItem creates a new order item value object
// This is the only way to create an OrderItem (factory method)
func NewOrderItem(productID string, quantity int, unitPrice float64) (OrderItem, error) {
	if quantity <= 0 {
		return OrderItem{}, errors.New("quantity must be positive")
	}
	if unitPrice < 0 {
		return OrderItem{}, errors.New("unit price cannot be negative")
	}

	return OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		SubTotal:  float64(quantity) * unitPrice,
	}, nil
}

// Validate validates the order item value object
func (item OrderItem) Validate() error {
	if item.ProductID == "" {
		return errors.New("product ID cannot be empty")
	}
	if item.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if item.UnitPrice < 0 {
		return errors.New("unit price cannot be negative")
	}
	if item.SubTotal != float64(item.Quantity)*item.UnitPrice {
		return errors.New("subtotal calculation mismatch")
	}
	return nil
}

