package mappers

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
)

// OrderDTOMapper defines the contract for mapping between domain Order and gRPC DTOs
type OrderDTOMapper interface {
	ToOrder(o order.Order) *ecommercev1.Order
	ToOrderList(orders []order.Order) []*ecommercev1.Order
	ToCreateOrderResponse(o order.Order) *ecommercev1.CreateOrderResponse
	ToGetOrderResponse(o order.Order) *ecommercev1.GetOrderResponse
	ToListUserOrdersResponse(orders []order.Order) *ecommercev1.ListUserOrdersResponse
	ToAddOrderItemResponse(o order.Order) *ecommercev1.AddOrderItemResponse
	ToConfirmOrderResponse(o order.Order) *ecommercev1.ConfirmOrderResponse
	ToCancelOrderResponse(o order.Order) *ecommercev1.CancelOrderResponse
}

