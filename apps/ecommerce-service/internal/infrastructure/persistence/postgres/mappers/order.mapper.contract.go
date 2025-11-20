package mappers

import "golang-social-media/apps/ecommerce-service/internal/domain/order"

// OrderMapper defines the contract for mapping between domain Order and persistence models
type OrderMapper interface {
	ToModel(o order.Order) (OrderModel, []OrderItemModel)
	ToDomain(orderModel OrderModel, itemModels []OrderItemModel) order.Order
	ToDomainList(orders []OrderModel, itemsMap map[string][]OrderItemModel) []order.Order
}


