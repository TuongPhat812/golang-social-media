package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"gorm.io/gorm"
)

// BatchOperations provides batch operations for repositories
type BatchOperations struct {
	db *gorm.DB
}

// NewBatchOperations creates a new BatchOperations instance
func NewBatchOperations(db *gorm.DB) *BatchOperations {
	return &BatchOperations{db: db}
}

// BatchCreateProducts creates multiple products in a single transaction
func (b *BatchOperations) BatchCreateProducts(ctx context.Context, products []product.Product, batchSize int) error {
	if len(products) == 0 {
		return nil
	}

	// Convert to models
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

	// Process in batches
	return b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(models); i += batchSize {
			end := i + batchSize
			if end > len(models) {
				end = len(models)
			}
			batch := models[i:end]
			if err := tx.Create(&batch).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdateProducts updates multiple products in a single transaction
func (b *BatchOperations) BatchUpdateProducts(ctx context.Context, products []product.Product, batchSize int) error {
	if len(products) == 0 {
		return nil
	}

	// Convert to models
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

	// Process in batches
	return b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(models); i += batchSize {
			end := i + batchSize
			if end > len(models) {
				end = len(models)
			}
			batch := models[i:end]
			if err := tx.Save(&batch).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdateProductStock updates stock for multiple products
func (b *BatchOperations) BatchUpdateProductStock(ctx context.Context, updates map[string]int) error {
	if len(updates) == 0 {
		return nil
	}

	return b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for productID, newStock := range updates {
			if err := tx.Model(&ProductModel{}).
				Where("id = ?", productID).
				Update("stock", newStock).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

