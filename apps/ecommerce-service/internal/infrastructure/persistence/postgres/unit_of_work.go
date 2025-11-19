package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/application/unit_of_work"

	"gorm.io/gorm"
)

var _ unit_of_work.UnitOfWork = (*unitOfWork)(nil)
var _ unit_of_work.Factory = (*UnitOfWorkFactory)(nil)

// unitOfWork implements UnitOfWork interface
type unitOfWork struct {
	db          *gorm.DB
	tx          *gorm.DB
	productRepo products.Repository
	orderRepo   orders.Repository
	committed   bool
	rolledBack  bool
}

// UnitOfWorkFactory creates new UnitOfWork instances
type UnitOfWorkFactory struct {
	db *gorm.DB
}

// NewUnitOfWorkFactory creates a new UnitOfWorkFactory
func NewUnitOfWorkFactory(db *gorm.DB) *UnitOfWorkFactory {
	return &UnitOfWorkFactory{db: db}
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

	// Create repositories with transaction (cache is not used in transactions)
	uow.productRepo = NewProductRepositoryWithTx(tx, nil)
	uow.orderRepo = NewOrderRepositoryWithTx(tx, nil)

	return uow, nil
}

// Products returns the product repository within this unit of work
func (u *unitOfWork) Products() products.Repository {
	return u.productRepo
}

// Orders returns the order repository within this unit of work
func (u *unitOfWork) Orders() orders.Repository {
	return u.orderRepo
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
