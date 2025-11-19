package factories

import (
	"time"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"github.com/google/uuid"
)

// OrderFactory creates Order aggregates with proper initialization
type OrderFactory struct{}

// NewOrderFactory creates a new OrderFactory
func NewOrderFactory() *OrderFactory {
	return &OrderFactory{}
}

// CreateOrder creates a new Order with proper initialization
// This factory encapsulates the complex creation logic
func (f *OrderFactory) CreateOrder(userID string) (*order.Order, error) {
	if userID == "" {
		return nil, &OrderFactoryError{Message: "user ID cannot be empty"}
	}

	now := time.Now().UTC()
	orderModel := &order.Order{
		ID:          uuid.NewString(),
		UserID:      userID,
		Status:      order.StatusDraft,
		Items:       []order.OrderItem{},
		TotalAmount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Validate the created order
	if err := orderModel.Validate(); err != nil {
		return nil, &OrderFactoryError{
			Message: "failed to validate order",
			Cause:   err,
		}
	}

	// Domain logic: create order (this adds domain events internally)
	orderModel.Create()

	return orderModel, nil
}

// OrderFactoryError represents an error in order factory
type OrderFactoryError struct {
	Message string
	Cause   error
}

func (e *OrderFactoryError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *OrderFactoryError) Unwrap() error {
	return e.Cause
}

