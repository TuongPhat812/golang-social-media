package postgres

import (
	"context"

	"github.com/google/uuid"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/cache"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers"
	"gorm.io/gorm"
)

var _ orders.Repository = (*OrderRepository)(nil)

type OrderRepository struct {
	db     *gorm.DB
	mapper *mappers.OrderMapper
	cache  *cache.OrderCache
}

func NewOrderRepository(db *gorm.DB, orderCache *cache.OrderCache) *OrderRepository {
	return &OrderRepository{
		db:     db,
		mapper: mappers.NewOrderMapper(),
		cache:  orderCache,
	}
}

// NewOrderRepositoryWithTx creates an OrderRepository with a specific transaction
func NewOrderRepositoryWithTx(tx *gorm.DB, orderCache *cache.OrderCache) *OrderRepository {
	return &OrderRepository{
		db:     tx,
		mapper: mappers.NewOrderMapper(),
		cache:  orderCache,
	}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	mapperOrderModel, mapperItemModels := r.mapper.ToModel(*o)

	// Convert to database models
	orderModel := OrderModel{
		ID:          mapperOrderModel.ID,
		UserID:      mapperOrderModel.UserID,
		Status:      mapperOrderModel.Status,
		TotalAmount: mapperOrderModel.TotalAmount,
		CreatedAt:   mapperOrderModel.CreatedAt,
		UpdatedAt:   mapperOrderModel.UpdatedAt,
	}
	itemModels := make([]OrderItemModel, len(mapperItemModels))
	for i, m := range mapperItemModels {
		itemModels[i] = OrderItemModel{
			ID:        m.ID,
			OrderID:   m.OrderID,
			ProductID: m.ProductID,
			Quantity:  m.Quantity,
			UnitPrice: m.UnitPrice,
			SubTotal:  m.SubTotal,
			CreatedAt: m.CreatedAt,
		}
	}

	// Generate IDs for order items
	for i := range itemModels {
		if itemModels[i].ID == "" {
			itemModels[i].ID = uuid.NewString()
		}
	}

	// Create order (transaction managed by UnitOfWork)
	if err := r.db.WithContext(ctx).Create(&orderModel).Error; err != nil {
		return err
	}

	// Create order items
	if len(itemModels) > 0 {
		if err := r.db.WithContext(ctx).Create(&itemModels).Error; err != nil {
			return err
		}
	}

	// Convert back to mapper models for ToDomain
	mapperOrderModel.ID = orderModel.ID
	mapperItemModels = make([]mappers.OrderItemModel, len(itemModels))
	for i, item := range itemModels {
		mapperItemModels[i] = mappers.OrderItemModel{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			SubTotal:  item.SubTotal,
			CreatedAt: item.CreatedAt,
		}
	}

	// Update domain entity
	*o = r.mapper.ToDomain(mapperOrderModel, mapperItemModels)
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (order.Order, error) {
	// Try cache first
	if r.cache != nil {
		if cached, err := r.cache.GetOrder(ctx, id); err == nil {
			return *cached, nil
		}
	}

	// Fallback to database
	var orderModel OrderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&orderModel).Error; err != nil {
		return order.Order{}, err
	}

	var itemModels []OrderItemModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).Find(&itemModels).Error; err != nil {
		return order.Order{}, err
	}

	// Convert to mapper models
	mapperOrderModel := mappers.OrderModel{
		ID:          orderModel.ID,
		UserID:      orderModel.UserID,
		Status:      orderModel.Status,
		TotalAmount: orderModel.TotalAmount,
		CreatedAt:   orderModel.CreatedAt,
		UpdatedAt:   orderModel.UpdatedAt,
	}
	mapperItemModels := make([]mappers.OrderItemModel, len(itemModels))
	for i, item := range itemModels {
		mapperItemModels[i] = mappers.OrderItemModel{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			SubTotal:  item.SubTotal,
			CreatedAt: item.CreatedAt,
		}
	}

	domainOrder := r.mapper.ToDomain(mapperOrderModel, mapperItemModels)

	// Cache the result
	if r.cache != nil {
		_ = r.cache.SetOrder(ctx, &domainOrder)
	}

	return domainOrder, nil
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	mapperOrderModel, mapperItemModels := r.mapper.ToModel(*o)

	// Convert to database models
	orderModel := OrderModel{
		ID:          mapperOrderModel.ID,
		UserID:      mapperOrderModel.UserID,
		Status:      mapperOrderModel.Status,
		TotalAmount: mapperOrderModel.TotalAmount,
		CreatedAt:   mapperOrderModel.CreatedAt,
		UpdatedAt:   mapperOrderModel.UpdatedAt,
	}
	itemModels := make([]OrderItemModel, len(mapperItemModels))
	for i, m := range mapperItemModels {
		itemModels[i] = OrderItemModel{
			ID:        m.ID,
			OrderID:   m.OrderID,
			ProductID: m.ProductID,
			Quantity:  m.Quantity,
			UnitPrice: m.UnitPrice,
			SubTotal:  m.SubTotal,
			CreatedAt: m.CreatedAt,
		}
	}

	// Generate IDs for new order items
	for i := range itemModels {
		if itemModels[i].ID == "" {
			itemModels[i].ID = uuid.NewString()
		}
	}

	// Update order (transaction managed by UnitOfWork)
	if err := r.db.WithContext(ctx).Save(&orderModel).Error; err != nil {
		return err
	}

	// Delete existing items
	if err := r.db.WithContext(ctx).Where("order_id = ?", o.ID).Delete(&OrderItemModel{}).Error; err != nil {
		return err
	}

	// Create new items
	if len(itemModels) > 0 {
		if err := r.db.WithContext(ctx).Create(&itemModels).Error; err != nil {
			return err
		}
	}

	// Convert back to mapper models for ToDomain
	mapperItemModels = make([]mappers.OrderItemModel, len(itemModels))
	for i, item := range itemModels {
		mapperItemModels[i] = mappers.OrderItemModel{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			SubTotal:  item.SubTotal,
			CreatedAt: item.CreatedAt,
		}
	}

	// Update domain entity
	*o = r.mapper.ToDomain(mapperOrderModel, mapperItemModels)

	// Invalidate cache
	if r.cache != nil {
		_ = r.cache.DeleteOrder(ctx, o.ID)
		_ = r.cache.InvalidateOrderList(ctx, o.UserID)
		_ = r.cache.SetOrder(ctx, o)
	}

	return nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string, limit int) ([]order.Order, error) {
	// Try cache first
	if r.cache != nil {
		if cached, err := r.cache.GetOrderList(ctx, userID, limit); err == nil {
			return cached, nil
		}
	}

	// Fallback to database
	var orderModels []OrderModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&orderModels).Error; err != nil {
		return nil, err
	}

	// Build items map for efficient lookup
	itemsMap := make(map[string][]OrderItemModel)
	for _, orderModel := range orderModels {
		var itemModels []OrderItemModel
		if err := r.db.WithContext(ctx).Where("order_id = ?", orderModel.ID).Find(&itemModels).Error; err != nil {
			return nil, err
		}
		itemsMap[orderModel.ID] = itemModels
	}

	// Convert to mapper models
	mapperOrderModels := make([]mappers.OrderModel, len(orderModels))
	for i, om := range orderModels {
		mapperOrderModels[i] = mappers.OrderModel{
			ID:          om.ID,
			UserID:      om.UserID,
			Status:      om.Status,
			TotalAmount: om.TotalAmount,
			CreatedAt:   om.CreatedAt,
			UpdatedAt:   om.UpdatedAt,
		}
	}
	mapperItemsMap := make(map[string][]mappers.OrderItemModel)
	for orderID, items := range itemsMap {
		mapperItems := make([]mappers.OrderItemModel, len(items))
		for i, item := range items {
			mapperItems[i] = mappers.OrderItemModel{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				SubTotal:  item.SubTotal,
				CreatedAt: item.CreatedAt,
			}
		}
		mapperItemsMap[orderID] = mapperItems
	}

	orders := r.mapper.ToDomainList(mapperOrderModels, mapperItemsMap)

	// Cache the result
	if r.cache != nil {
		_ = r.cache.SetOrderList(ctx, userID, limit, orders)
	}

	return orders, nil
}

