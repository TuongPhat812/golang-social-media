package ecommerce

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/logger"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	ecommercev1.UnimplementedOrderServiceServer
	deps *bootstrap.Dependencies
	log  *zerolog.Logger
}

func NewOrderHandler(deps *bootstrap.Dependencies) *OrderHandler {
	return &OrderHandler{
		deps: deps,
		log:  logger.Component("ecommerce.grpc.order"),
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *ecommercev1.CreateOrderRequest) (*ecommercev1.CreateOrderResponse, error) {
	cmdReq := contracts.CreateOrderCommandRequest{
		UserID: req.GetUserId(),
	}

	order, err := h.deps.CreateOrderCmd.Execute(ctx, cmdReq)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("user_id", req.GetUserId()).
			Msg("failed to create order")
		return nil, err
	}

	h.log.Info().
		Str("order_id", order.ID).
		Str("user_id", order.UserID).
		Msg("order created")

	return &ecommercev1.CreateOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *ecommercev1.GetOrderRequest) (*ecommercev1.GetOrderResponse, error) {
	order, err := h.deps.GetOrderQuery.Execute(ctx, req.GetOrderId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Msg("failed to get order")
		return nil, err
	}

	return &ecommercev1.GetOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

func (h *OrderHandler) ListUserOrders(ctx context.Context, req *ecommercev1.ListUserOrdersRequest) (*ecommercev1.ListUserOrdersResponse, error) {
	orders, err := h.deps.ListUserOrdersQuery.Execute(ctx, req.GetUserId(), int(req.GetLimit()))
	if err != nil {
		h.log.Error().
			Err(err).
			Str("user_id", req.GetUserId()).
			Msg("failed to list user orders")
		return nil, err
	}

	pbOrders := make([]*ecommercev1.Order, len(orders))
	for i, o := range orders {
		pbOrders[i] = h.orderToProto(o)
	}

	return &ecommercev1.ListUserOrdersResponse{
		Orders: pbOrders,
	}, nil
}

func (h *OrderHandler) AddOrderItem(ctx context.Context, req *ecommercev1.AddOrderItemRequest) (*ecommercev1.AddOrderItemResponse, error) {
	cmdReq := contracts.AddOrderItemCommandRequest{
		OrderID:   req.GetOrderId(),
		ProductID: req.GetProductId(),
		Quantity:  int(req.GetQuantity()),
	}

	order, err := h.deps.AddOrderItemCmd.Execute(ctx, cmdReq)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Str("product_id", req.GetProductId()).
			Msg("failed to add order item")
		return nil, err
	}

	h.log.Info().
		Str("order_id", order.ID).
		Str("product_id", req.GetProductId()).
		Int("quantity", int(req.GetQuantity())).
		Msg("order item added")

	return &ecommercev1.AddOrderItemResponse{
		Order: h.orderToProto(order),
	}, nil
}

func (h *OrderHandler) ConfirmOrder(ctx context.Context, req *ecommercev1.ConfirmOrderRequest) (*ecommercev1.ConfirmOrderResponse, error) {
	err := h.deps.ConfirmOrderCmd.Execute(ctx, req.GetOrderId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Msg("failed to confirm order")
		return nil, err
	}

	// Fetch updated order
	order, err := h.deps.GetOrderQuery.Execute(ctx, req.GetOrderId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Msg("failed to get updated order")
		return nil, err
	}

	h.log.Info().
		Str("order_id", order.ID).
		Msg("order confirmed")

	return &ecommercev1.ConfirmOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *ecommercev1.CancelOrderRequest) (*ecommercev1.CancelOrderResponse, error) {
	err := h.deps.CancelOrderCmd.Execute(ctx, req.GetOrderId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Msg("failed to cancel order")
		return nil, err
	}

	// Fetch updated order
	order, err := h.deps.GetOrderQuery.Execute(ctx, req.GetOrderId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", req.GetOrderId()).
			Msg("failed to get updated order")
		return nil, err
	}

	h.log.Info().
		Str("order_id", order.ID).
		Msg("order cancelled")

	return &ecommercev1.CancelOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

func (h *OrderHandler) orderToProto(o order.Order) *ecommercev1.Order {
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

