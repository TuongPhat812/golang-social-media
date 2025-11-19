package postgres

import (
	"context"

	"github.com/google/uuid"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"gorm.io/gorm"
)

var _ orders.Repository = (*OrderRepository)(nil)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	orderModel, itemModels := OrderModelFromDomain(*o)

	// Generate IDs for order items
	for i := range itemModels {
		itemModels[i].ID = uuid.NewString()
	}

	// Use transaction to ensure consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(&orderModel).Error; err != nil {
			return err
		}

		// Create order items
		if len(itemModels) > 0 {
			if err := tx.Create(&itemModels).Error; err != nil {
				return err
			}
		}

		// Update domain entity
		*o = orderModel.ToDomain(itemModels)
		return nil
	})
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (order.Order, error) {
	var orderModel OrderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&orderModel).Error; err != nil {
		return order.Order{}, err
	}

	var itemModels []OrderItemModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).Find(&itemModels).Error; err != nil {
		return order.Order{}, err
	}

	return orderModel.ToDomain(itemModels), nil
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	orderModel, itemModels := OrderModelFromDomain(*o)

	// Generate IDs for new order items
	for i := range itemModels {
		if itemModels[i].ID == "" {
			itemModels[i].ID = uuid.NewString()
		}
	}

	// Use transaction to ensure consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order
		if err := tx.Save(&orderModel).Error; err != nil {
			return err
		}

		// Delete existing items
		if err := tx.Where("order_id = ?", o.ID).Delete(&OrderItemModel{}).Error; err != nil {
			return err
		}

		// Create new items
		if len(itemModels) > 0 {
			if err := tx.Create(&itemModels).Error; err != nil {
				return err
			}
		}

		// Update domain entity
		*o = orderModel.ToDomain(itemModels)
		return nil
	})
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string, limit int) ([]order.Order, error) {
	var orderModels []OrderModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&orderModels).Error; err != nil {
		return nil, err
	}

	orders := make([]order.Order, len(orderModels))
	for i, orderModel := range orderModels {
		var itemModels []OrderItemModel
		if err := r.db.WithContext(ctx).Where("order_id = ?", orderModel.ID).Find(&itemModels).Error; err != nil {
			return nil, err
		}
		orders[i] = orderModel.ToDomain(itemModels)
	}

	return orders, nil
}

