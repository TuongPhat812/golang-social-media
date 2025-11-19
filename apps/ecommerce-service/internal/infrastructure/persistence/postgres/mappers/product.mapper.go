package mappers

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
)

// ProductMapper maps between domain Product and persistence models
type ProductMapper struct{}

// NewProductMapper creates a new ProductMapper
func NewProductMapper() *ProductMapper {
	return &ProductMapper{}
}

// ToModel converts a domain Product to ProductModel
func (m *ProductMapper) ToModel(p product.Product) postgres.ProductModel {
	return postgres.ProductModel{
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
func (m *ProductMapper) ToDomain(model postgres.ProductModel) product.Product {
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
func (m *ProductMapper) ToDomainList(models []postgres.ProductModel) []product.Product {
	products := make([]product.Product, len(models))
	for i, model := range models {
		products[i] = m.ToDomain(model)
	}
	return products
}

// ToModelList converts a slice of domain Products to ProductModels
func (m *ProductMapper) ToModelList(products []product.Product) []postgres.ProductModel {
	models := make([]postgres.ProductModel, len(products))
	for i, p := range products {
		models[i] = m.ToModel(p)
	}
	return models
}

