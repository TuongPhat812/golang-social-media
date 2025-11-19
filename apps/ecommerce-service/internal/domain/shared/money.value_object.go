package shared

import (
	"errors"
	"fmt"
	"strings"
)

// Money represents a monetary value with currency (Value Object)
// Value Objects are immutable and defined by their values, not identity
type Money struct {
	amount   float64
	currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	// Normalize currency to uppercase
	normalizedCurrency := strings.ToUpper(currency)
	return Money{
		amount:   amount,
		currency: normalizedCurrency,
	}, nil
}

// Amount returns the amount
func (m Money) Amount() float64 {
	return m.amount
}

// Currency returns the currency
func (m Money) Currency() string {
	return m.currency
}

// String returns a string representation of Money
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Add adds two Money values (must have same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract subtracts two Money values (must have same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	if m.amount < other.amount {
		return Money{}, errors.New("insufficient money")
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

// Multiply multiplies Money by a scalar
func (m Money) Multiply(scalar float64) (Money, error) {
	if scalar < 0 {
		return Money{}, errors.New("scalar cannot be negative")
	}
	return NewMoney(m.amount*scalar, m.currency)
}

// IsZero returns true if amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// Compare compares two Money values
// Returns: -1 if m < other, 0 if m == other, 1 if m > other
func (m Money) Compare(other Money) (int, error) {
	if m.currency != other.currency {
		return 0, errors.New("cannot compare money with different currencies")
	}
	if m.amount < other.amount {
		return -1, nil
	}
	if m.amount > other.amount {
		return 1, nil
	}
	return 0, nil
}

