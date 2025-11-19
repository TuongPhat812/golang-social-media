package services

import (
	"errors"
)

// StockReservationService handles stock reservation logic
// This is a domain service because it coordinates between multiple entities
type StockReservationService struct{}

// NewStockReservationService creates a new stock reservation service
func NewStockReservationService() *StockReservationService {
	return &StockReservationService{}
}

// ReserveStock validates if stock can be reserved for an order
// This logic doesn't belong to Product or Order entity alone
func (s *StockReservationService) ReserveStock(
	availableStock int,
	requestedQuantity int,
) error {
	if availableStock < requestedQuantity {
		return errors.New("insufficient stock available")
	}
	if requestedQuantity <= 0 {
		return errors.New("requested quantity must be positive")
	}
	return nil
}

// CalculateReservedStock calculates total reserved stock for multiple items
func (s *StockReservationService) CalculateReservedStock(
	itemQuantities map[string]int, // productID -> quantity
) int {
	total := 0
	for _, qty := range itemQuantities {
		total += qty
	}
	return total
}

