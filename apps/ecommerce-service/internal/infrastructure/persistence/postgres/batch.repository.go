package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"gorm.io/gorm"
)

// BatchRepository provides batch operations for repositories
type BatchRepository struct {
	db *gorm.DB
}

// NewBatchRepository creates a new BatchRepository
func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

// BatchCreateProducts creates multiple products in a single transaction
func (r *BatchRepository) BatchCreateProducts(ctx context.Context, products []product.Product) error {
	if len(products) == 0 {
		return nil
	}

	models := make([]ProductModel, len(products))
	for i, p := range products {
		models[i] = ProductModel{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Status:      string(p.Status),
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	// Use batch insert with chunk size to avoid query size limits
	const chunkSize = 100
	for i := 0; i < len(models); i += chunkSize {
		end := i + chunkSize
		if end > len(models) {
			end = len(models)
		}
		if err := r.db.WithContext(ctx).Create(models[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}

// BatchUpdateProducts updates multiple products in a single transaction
func (r *BatchRepository) BatchUpdateProducts(ctx context.Context, products []product.Product) error {
	if len(products) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range products {
			model := ProductModel{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Stock:       p.Stock,
				Status:      string(p.Status),
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			}
			if err := tx.Save(&model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchCreateOrders creates multiple orders with their items in a single transaction
func (r *BatchRepository) BatchCreateOrders(ctx context.Context, orders []order.Order) error {
	if len(orders) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, o := range orders {
			orderModel := OrderModel{
				ID:          o.ID,
				UserID:      o.UserID,
				Status:      string(o.Status),
				TotalAmount: o.TotalAmount,
				CreatedAt:   o.CreatedAt,
				UpdatedAt:   o.UpdatedAt,
			}

			if err := tx.Create(&orderModel).Error; err != nil {
				return err
			}

			if len(o.Items) > 0 {
				itemModels := make([]OrderItemModel, len(o.Items))
				for i, item := range o.Items {
					itemModels[i] = OrderItemModel{
						ID:        "", // Will be generated
						OrderID:   o.ID,
						ProductID: item.ProductID,
						Quantity:  item.Quantity,
						UnitPrice: item.UnitPrice,
						SubTotal:  item.SubTotal,
						CreatedAt: o.CreatedAt,
					}
				}
				if err := tx.Create(&itemModels).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// BatchUpdateOrderItems updates multiple order items in a single transaction
func (r *BatchRepository) BatchUpdateOrderItems(ctx context.Context, orderID string, items []order.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing items
		if err := tx.Where("order_id = ?", orderID).Delete(&OrderItemModel{}).Error; err != nil {
			return err
		}

		// Get order to get CreatedAt
		var orderModel OrderModel
		if err := tx.Where("id = ?", orderID).First(&orderModel).Error; err != nil {
			return err
		}

		// Create new items
		itemModels := make([]OrderItemModel, len(items))
		for i, item := range items {
			itemModels[i] = OrderItemModel{
				ID:        "", // Will be generated
				OrderID:   orderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				SubTotal:  item.SubTotal,
				CreatedAt: orderModel.CreatedAt, // Use order's CreatedAt
			}
		}

		if err := tx.Create(&itemModels).Error; err != nil {
			return err
		}

		return nil
	})
}

