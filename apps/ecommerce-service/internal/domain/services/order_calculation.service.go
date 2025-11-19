package services

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/apps/ecommerce-service/internal/domain/shared"
)

// OrderCalculationService handles order calculation logic
// This is a domain service because it involves complex business rules
type OrderCalculationService struct {
	defaultCurrency string
}

// NewOrderCalculationService creates a new order calculation service
func NewOrderCalculationService(defaultCurrency string) *OrderCalculationService {
	return &OrderCalculationService{
		defaultCurrency: defaultCurrency,
	}
}

// CalculateOrderTotal calculates the total amount for an order
// This could include discounts, taxes, shipping, etc.
func (s *OrderCalculationService) CalculateOrderTotal(items []order.OrderItem) (shared.Money, error) {
	total, err := shared.NewMoney(0, s.defaultCurrency)
	if err != nil {
		return shared.Money{}, err
	}

	for _, item := range items {
		itemPrice, err := shared.NewMoney(item.UnitPrice, s.defaultCurrency)
		if err != nil {
			return shared.Money{}, err
		}

		quantity, err := shared.NewQuantity(item.Quantity)
		if err != nil {
			return shared.Money{}, err
		}

		// Calculate subtotal: price * quantity
		multiplied, err := itemPrice.Multiply(float64(quantity.Value()))
		if err != nil {
			return shared.Money{}, err
		}

		total, err = total.Add(multiplied)
		if err != nil {
			return shared.Money{}, err
		}
	}

	return total, nil
}

// CalculateItemSubtotal calculates subtotal for a single order item
func (s *OrderCalculationService) CalculateItemSubtotal(
	unitPrice float64,
	quantity int,
) (shared.Money, error) {
	price, err := shared.NewMoney(unitPrice, s.defaultCurrency)
	if err != nil {
		return shared.Money{}, err
	}

	qty, err := shared.NewQuantity(quantity)
	if err != nil {
		return shared.Money{}, err
	}

	return price.Multiply(float64(qty.Value()))
}

