package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers"
	"gorm.io/gorm"
)

var _ products.Repository = (*ProductRepository)(nil)

type ProductRepository struct {
	db     *gorm.DB
	mapper *mappers.ProductMapper
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db:     db,
		mapper: mappers.NewProductMapper(),
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	model := r.mapper.ToModel(*p)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	*p = r.mapper.ToDomain(model)
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (product.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return product.Product{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	model := r.mapper.ToModel(*p)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}
	*p = r.mapper.ToDomain(model)
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

	return r.mapper.ToDomainList(models), nil
}

