package postgres

import (
	"time"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

type ProductModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;type:text;not null"`
	Description string    `gorm:"column:description;type:text"`
	Price       float64   `gorm:"column:price;type:decimal(10,2);not null"`
	Stock       int       `gorm:"column:stock;type:integer;not null;default:0"`
	Status      string    `gorm:"column:status;type:text;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (ProductModel) TableName() string {
	return "products"
}

func ProductModelFromDomain(p product.Product) ProductModel {
	return ProductModel{
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

func (m ProductModel) ToDomain() product.Product {
	return product.Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		Stock:       m.Stock,
		Status:      product.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

