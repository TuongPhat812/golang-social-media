package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"gorm.io/gorm"
)

var _ products.Repository = (*ProductRepository)(nil)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	model := ProductModelFromDomain(*p)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	*p = model.ToDomain()
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (product.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return product.Product{}, err
	}
	return model.ToDomain(), nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	model := ProductModelFromDomain(*p)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}
	*p = model.ToDomain()
	return nil
}

func (r *ProductRepository) List(ctx context.Context, status *product.Status, limit, offset int) ([]product.Product, error) {
	var models []ProductModel
	query := r.db.WithContext(ctx)

	if status != nil {
		query = query.Where("status = ?", string(*status))
	}

	if err := query.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]product.Product, len(models))
	for i, model := range models {
		products[i] = model.ToDomain()
	}

	return products, nil
}

