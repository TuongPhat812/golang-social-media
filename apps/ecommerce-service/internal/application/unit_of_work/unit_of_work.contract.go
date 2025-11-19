package unit_of_work

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
)

// UnitOfWork manages a transaction and provides access to repositories
// All repositories returned from a UnitOfWork share the same transaction
type UnitOfWork interface {
	// Products returns the product repository within this unit of work
	Products() products.Repository

	// Orders returns the order repository within this unit of work
	Orders() orders.Repository

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}

// Factory creates new UnitOfWork instances
type Factory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

