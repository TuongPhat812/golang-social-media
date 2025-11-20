package mappers

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrderDTOMapperImpl implements OrderDTOMapper interface
type OrderDTOMapperImpl struct{}

var _ OrderDTOMapper = (*OrderDTOMapperImpl)(nil)

// NewOrderDTOMapper creates a new OrderDTOMapperImpl
func NewOrderDTOMapper() OrderDTOMapper {
	return &OrderDTOMapperImpl{}
}

// FromCreateOrderRequest extracts user ID from CreateOrderRequest
func (m *OrderDTOMapperImpl) FromCreateOrderRequest(req *ecommercev1.CreateOrderRequest) string {
	return req.GetUserId()
}

// FromAddOrderItemRequest converts gRPC AddOrderItemRequest to command request data
func (m *OrderDTOMapperImpl) FromAddOrderItemRequest(req *ecommercev1.AddOrderItemRequest) (orderID, productID string, quantity int) {
	return req.GetOrderId(), req.GetProductId(), int(req.GetQuantity())
}

// ToOrder converts domain Order to gRPC Order
func (m *OrderDTOMapperImpl) ToOrder(o order.Order) *ecommercev1.Order {
	pbItems := make([]*ecommercev1.OrderItem, len(o.Items))
	for i, item := range o.Items {
		pbItems[i] = &ecommercev1.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			UnitPrice: item.UnitPrice,
			SubTotal:  item.SubTotal,
		}
	}

	return &ecommercev1.Order{
		Id:          o.ID,
		UserId:      o.UserID,
		Status:      string(o.Status),
		Items:       pbItems,
		TotalAmount: o.TotalAmount,
		CreatedAt:   timestamppb.New(o.CreatedAt),
		UpdatedAt:   timestamppb.New(o.UpdatedAt),
	}
}

// ToOrderList converts a slice of domain Orders to gRPC Orders
func (m *OrderDTOMapperImpl) ToOrderList(orders []order.Order) []*ecommercev1.Order {
	result := make([]*ecommercev1.Order, len(orders))
	for i, o := range orders {
		result[i] = m.ToOrder(o)
	}
	return result
}

// ToCreateOrderResponse converts domain Order to CreateOrderResponse
func (m *OrderDTOMapperImpl) ToCreateOrderResponse(o order.Order) *ecommercev1.CreateOrderResponse {
	return &ecommercev1.CreateOrderResponse{
		Order: m.ToOrder(o),
	}
}

// ToGetOrderResponse converts domain Order to GetOrderResponse
func (m *OrderDTOMapperImpl) ToGetOrderResponse(o order.Order) *ecommercev1.GetOrderResponse {
	return &ecommercev1.GetOrderResponse{
		Order: m.ToOrder(o),
	}
}

// ToListUserOrdersResponse converts domain Orders to ListUserOrdersResponse
func (m *OrderDTOMapperImpl) ToListUserOrdersResponse(orders []order.Order) *ecommercev1.ListUserOrdersResponse {
	return &ecommercev1.ListUserOrdersResponse{
		Orders: m.ToOrderList(orders),
	}
}

// ToAddOrderItemResponse converts domain Order to AddOrderItemResponse
func (m *OrderDTOMapperImpl) ToAddOrderItemResponse(o order.Order) *ecommercev1.AddOrderItemResponse {
	return &ecommercev1.AddOrderItemResponse{
		Order: m.ToOrder(o),
	}
}

// ToConfirmOrderResponse converts domain Order to ConfirmOrderResponse
func (m *OrderDTOMapperImpl) ToConfirmOrderResponse(o order.Order) *ecommercev1.ConfirmOrderResponse {
	return &ecommercev1.ConfirmOrderResponse{
		Order: m.ToOrder(o),
	}
}

// ToCancelOrderResponse converts domain Order to CancelOrderResponse
func (m *OrderDTOMapperImpl) ToCancelOrderResponse(o order.Order) *ecommercev1.CancelOrderResponse {
	return &ecommercev1.CancelOrderResponse{
		Order: m.ToOrder(o),
	}
}

