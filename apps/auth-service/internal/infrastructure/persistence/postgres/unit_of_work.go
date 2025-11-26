package postgres

import (
	"context"
	"encoding/json"

	"golang-social-media/apps/auth-service/internal/application/repository"
	"golang-social-media/apps/auth-service/internal/application/unit_of_work"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres/mappers"

	"gorm.io/gorm"
)

var _ unit_of_work.UnitOfWork = (*unitOfWork)(nil)
var _ unit_of_work.Factory = (*UnitOfWorkFactory)(nil)

// unitOfWork implements UnitOfWork interface
type unitOfWork struct {
	db          *gorm.DB
	tx          *gorm.DB
	userRepo    repository.UserRepository
	outboxRepo  *OutboxRepository
	eventStoreRepo *EventStoreRepository
	committed   bool
	rolledBack  bool
}

// UnitOfWorkFactory creates new UnitOfWork instances
type UnitOfWorkFactory struct {
	db            *gorm.DB
	userMapper    mappers.UserMapper
}

// NewUnitOfWorkFactory creates a new UnitOfWorkFactory
func NewUnitOfWorkFactory(
	db *gorm.DB,
	userMapper mappers.UserMapper,
) *UnitOfWorkFactory {
	return &UnitOfWorkFactory{
		db:         db,
		userMapper: userMapper,
	}
}

// New creates a new UnitOfWork with a transaction
func (f *UnitOfWorkFactory) New(ctx context.Context) (unit_of_work.UnitOfWork, error) {
	tx := f.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	uow := &unitOfWork{
		db:         f.db,
		tx:         tx,
		committed:  false,
		rolledBack: false,
	}

	// Create repositories with transaction
	uow.userRepo = NewUserRepositoryWithTx(tx, f.userMapper, nil)
	uow.outboxRepo = NewOutboxRepositoryWithTx(tx)
	uow.eventStoreRepo = NewEventStoreRepositoryWithTx(tx)

	return uow, nil
}

// Users returns the user repository within this unit of work
func (u *unitOfWork) Users() repository.UserRepository {
	return u.userRepo
}

// SaveEvents saves domain events to outbox and event store within the transaction
func (u *unitOfWork) SaveEvents(ctx context.Context, events []interface{}) error {
	for _, event := range events {
		// Extract event information
		var aggregateID, aggregateType, eventType string
		var eventVersion int = 1

		// Try to get event type from interface
		if domainEvent, ok := event.(interface {
			Type() string
		}); ok {
			eventType = domainEvent.Type()
		}

		// Extract aggregate info from struct fields
		eventMap, err := structToMap(event)
		if err != nil {
			return err
		}

		// Try common field names for auth-service events
		if id, ok := eventMap["UserID"].(string); ok && id != "" {
			aggregateID = id
			aggregateType = "User"
		} else if id, ok := eventMap["RoleID"].(string); ok && id != "" {
			aggregateID = id
			aggregateType = "Role"
		} else if id, ok := eventMap["PermissionID"].(string); ok && id != "" {
			aggregateID = id
			aggregateType = "Permission"
		}

		// Save to outbox
		if err := u.outboxRepo.Create(ctx, aggregateID, aggregateType, eventType, eventVersion, event); err != nil {
			return err
		}

		// Save to event store
		metadata := map[string]interface{}{}
		if err := u.eventStoreRepo.Append(ctx, aggregateID, aggregateType, eventType, eventVersion, event, metadata); err != nil {
			return err
		}
	}

	return nil
}

// structToMap converts a struct to a map
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// Commit commits the transaction
func (u *unitOfWork) Commit() error {
	if u.committed {
		return nil // Already committed
	}
	if u.rolledBack {
		return nil // Already rolled back, nothing to commit
	}

	if err := u.tx.Commit().Error; err != nil {
		return err
	}

	u.committed = true
	return nil
}

// Rollback rolls back the transaction
func (u *unitOfWork) Rollback() error {
	if u.rolledBack {
		return nil // Already rolled back
	}
	if u.committed {
		return nil // Already committed, nothing to rollback
	}

	if err := u.tx.Rollback().Error; err != nil {
		return err
	}

	u.rolledBack = true
	return nil
}

