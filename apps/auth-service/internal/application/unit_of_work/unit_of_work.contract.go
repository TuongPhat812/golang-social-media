package unit_of_work

import (
	"context"

	"golang-social-media/apps/auth-service/internal/application/repository"
)

// UnitOfWork manages a transaction and provides access to repositories
// All repositories returned from a UnitOfWork share the same transaction
type UnitOfWork interface {
	// Users returns the user repository within this unit of work
	Users() repository.UserRepository

	// SaveEvents saves domain events to outbox and event store within the transaction
	SaveEvents(ctx context.Context, events []interface{}) error

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}

// Factory creates new UnitOfWork instances
type Factory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

