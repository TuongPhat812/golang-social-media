package services

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/apps/ecommerce-service/internal/domain/shared"
)

// ProductAvailabilityService handles product availability logic
// This is a domain service because it involves complex business rules
type ProductAvailabilityService struct{}

// NewProductAvailabilityService creates a new product availability service
func NewProductAvailabilityService() *ProductAvailabilityService {
	return &ProductAvailabilityService{}
}

// CheckAvailability checks if a product is available for purchase
// This could include checking stock, status, restrictions, etc.
func (s *ProductAvailabilityService) CheckAvailability(
	p product.Product,
	requestedQuantity shared.Quantity,
) (bool, string) {
	// Check if product is active
	if p.Status != product.StatusActive {
		return false, "product is not active"
	}

	// Check if product has stock
	if p.Stock == 0 {
		return false, "product is out of stock"
	}

	// Check if requested quantity is available
	availableQty, err := shared.NewQuantity(p.Stock)
	if err != nil {
		return false, "invalid stock quantity"
	}

	if availableQty.IsLessThan(requestedQuantity) {
		return false, "insufficient stock"
	}

	return true, ""
}

// CheckBulkAvailability checks availability for multiple products
func (s *ProductAvailabilityService) CheckBulkAvailability(
	products map[string]product.Product, // productID -> product
	requestedQuantities map[string]shared.Quantity, // productID -> quantity
) (map[string]bool, map[string]string) {
	available := make(map[string]bool)
	reasons := make(map[string]string)

	for productID, qty := range requestedQuantities {
		prod, exists := products[productID]
		if !exists {
			available[productID] = false
			reasons[productID] = "product not found"
			continue
		}

		isAvailable, reason := s.CheckAvailability(prod, qty)
		available[productID] = isAvailable
		if !isAvailable {
			reasons[productID] = reason
		}
	}

	return available, reasons
}

