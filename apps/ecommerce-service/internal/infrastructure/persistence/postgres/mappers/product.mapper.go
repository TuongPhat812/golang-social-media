package mappers

import (
	"time"

	"golang-social-media/apps/ecommerce-service/internal/domain/product"
)

// ProductModel represents the database model for Product
// This is defined here to avoid import cycle
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

// ProductMapper maps between domain Product and persistence models
type ProductMapper struct{}

// NewProductMapper creates a new ProductMapper
func NewProductMapper() *ProductMapper {
	return &ProductMapper{}
}

// ToModel converts a domain Product to ProductModel
func (m *ProductMapper) ToModel(p product.Product) ProductModel {
	return ProductModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price, // Currently using float64, can be migrated to Money value object
		Stock:       p.Stock, // Currently using int, can be migrated to Quantity value object
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToDomain converts a ProductModel to domain Product
func (m *ProductMapper) ToDomain(model ProductModel) product.Product {
	return product.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price, // Currently using float64
		Stock:       model.Stock, // Currently using int
		Status:      product.Status(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// ToDomainList converts a slice of ProductModel to domain Products
func (m *ProductMapper) ToDomainList(models []ProductModel) []product.Product {
	products := make([]product.Product, len(models))
	for i, model := range models {
		products[i] = m.ToDomain(model)
	}
	return products
}

// ToModelList converts a slice of domain Products to ProductModels
func (m *ProductMapper) ToModelList(products []product.Product) []ProductModel {
	models := make([]ProductModel, len(products))
	for i, p := range products {
		models[i] = m.ToModel(p)
	}
	return models
}

