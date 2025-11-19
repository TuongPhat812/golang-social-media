package shared

import (
	"errors"
	"fmt"
)

// Quantity represents a quantity value object
// Value Objects are immutable and defined by their values, not identity
type Quantity struct {
	value int
}

// NewQuantity creates a new Quantity value object
func NewQuantity(value int) (Quantity, error) {
	if value < 0 {
		return Quantity{}, errors.New("quantity cannot be negative")
	}
	return Quantity{value: value}, nil
}

// Value returns the quantity value
func (q Quantity) Value() int {
	return q.value
}

// String returns a string representation of Quantity
func (q Quantity) String() string {
	return fmt.Sprintf("%d", q.value)
}

// Add adds two quantities
func (q Quantity) Add(other Quantity) Quantity {
	return Quantity{value: q.value + other.value}
}

// Subtract subtracts two quantities
func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	if q.value < other.value {
		return Quantity{}, errors.New("insufficient quantity")
	}
	return Quantity{value: q.value - other.value}, nil
}

// Multiply multiplies quantity by a scalar
func (q Quantity) Multiply(scalar int) (Quantity, error) {
	if scalar < 0 {
		return Quantity{}, errors.New("scalar cannot be negative")
	}
	return Quantity{value: q.value * scalar}, nil
}

// IsZero returns true if quantity is zero
func (q Quantity) IsZero() bool {
	return q.value == 0
}

// Compare compares two quantities
// Returns: -1 if q < other, 0 if q == other, 1 if q > other
func (q Quantity) Compare(other Quantity) int {
	if q.value < other.value {
		return -1
	}
	if q.value > other.value {
		return 1
	}
	return 0
}

// IsGreaterThan returns true if quantity is greater than other
func (q Quantity) IsGreaterThan(other Quantity) bool {
	return q.Compare(other) > 0
}

// IsLessThan returns true if quantity is less than other
func (q Quantity) IsLessThan(other Quantity) bool {
	return q.Compare(other) < 0
}

