package postgres

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/cache"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers"
	"gorm.io/gorm"
)

var _ products.Repository = (*ProductRepository)(nil)

type ProductRepository struct {
	db     *gorm.DB
	mapper *mappers.ProductMapper
	cache  *cache.ProductCache
}

func NewProductRepository(db *gorm.DB, productCache *cache.ProductCache) *ProductRepository {
	return &ProductRepository{
		db:     db,
		mapper: mappers.NewProductMapper(),
		cache:  productCache,
	}
}

// NewProductRepositoryWithTx creates a ProductRepository with a specific transaction
func NewProductRepositoryWithTx(tx *gorm.DB, productCache *cache.ProductCache) *ProductRepository {
	return &ProductRepository{
		db:     tx,
		mapper: mappers.NewProductMapper(),
		cache:  productCache,
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	mapperModel := r.mapper.ToModel(*p)
	// Convert mappers.ProductModel to ProductModel
	model := ProductModel{
		ID:          mapperModel.ID,
		Name:        mapperModel.Name,
		Description: mapperModel.Description,
		Price:       mapperModel.Price,
		Stock:       mapperModel.Stock,
		Status:      mapperModel.Status,
		CreatedAt:   mapperModel.CreatedAt,
		UpdatedAt:   mapperModel.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	*p = r.mapper.ToDomain(mapperModel)

	// Invalidate cache
	if r.cache != nil {
		_ = r.cache.InvalidateProductList(ctx)
		_ = r.cache.SetProduct(ctx, p)
	}

	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (product.Product, error) {
	// Try cache first
	if r.cache != nil {
		if cached, err := r.cache.GetProduct(ctx, id); err == nil {
			return *cached, nil
		}
	}

	// Fallback to database
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return product.Product{}, err
	}
	// Convert ProductModel to mappers.ProductModel
	mapperModel := mappers.ProductModel{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Stock:       model.Stock,
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
	domainProduct := r.mapper.ToDomain(mapperModel)

	// Cache the result
	if r.cache != nil {
		_ = r.cache.SetProduct(ctx, &domainProduct)
	}

	return domainProduct, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	mapperModel := r.mapper.ToModel(*p)
	// Convert mappers.ProductModel to ProductModel
	model := ProductModel{
		ID:          mapperModel.ID,
		Name:        mapperModel.Name,
		Description: mapperModel.Description,
		Price:       mapperModel.Price,
		Stock:       mapperModel.Stock,
		Status:      mapperModel.Status,
		CreatedAt:   mapperModel.CreatedAt,
		UpdatedAt:   mapperModel.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}
	*p = r.mapper.ToDomain(mapperModel)

	// Invalidate cache
	if r.cache != nil {
		_ = r.cache.DeleteProduct(ctx, p.ID)
		_ = r.cache.InvalidateProductList(ctx)
		_ = r.cache.SetProduct(ctx, p)
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, status *product.Status, limit, offset int) ([]product.Product, error) {
	// Try cache first
	if r.cache != nil {
		if cached, err := r.cache.GetProductList(ctx, status, limit, offset); err == nil {
			return cached, nil
		}
	}

	// Fallback to database
	var models []ProductModel
	query := r.db.WithContext(ctx)

	if status != nil {
		query = query.Where("status = ?", string(*status))
	}

	if err := query.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}

	// Convert []ProductModel to []mappers.ProductModel
	mapperModels := make([]mappers.ProductModel, len(models))
	for i, m := range models {
		mapperModels[i] = mappers.ProductModel{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			Stock:       m.Stock,
			Status:      m.Status,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}

	products := r.mapper.ToDomainList(mapperModels)

	// Cache the result
	if r.cache != nil {
		_ = r.cache.SetProductList(ctx, status, limit, offset, products)
	}

	return products, nil
}

