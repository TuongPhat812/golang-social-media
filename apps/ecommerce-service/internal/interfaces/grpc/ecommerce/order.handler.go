package ecommerce

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	"golang-social-media/apps/ecommerce-service/internal/interfaces/grpc/mappers"
	"golang-social-media/pkg/logger"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
)

type OrderHandler struct {
	ecommercev1.UnimplementedOrderServiceServer
	deps      *bootstrap.Dependencies
	dtoMapper mappers.OrderDTOMapper
	log       *zerolog.Logger
}

func NewOrderHandler(deps *bootstrap.Dependencies, dtoMapper mappers.OrderDTOMapper) *OrderHandler {
	return &OrderHandler{
		deps:      deps,
		dtoMapper: dtoMapper,
		log:       logger.Component("ecommerce.grpc.order"),
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

	return h.dtoMapper.ToCreateOrderResponse(order), nil
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

	return h.dtoMapper.ToGetOrderResponse(order), nil
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

	return h.dtoMapper.ToListUserOrdersResponse(orders), nil
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

	return h.dtoMapper.ToAddOrderItemResponse(order), nil
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

	return h.dtoMapper.ToConfirmOrderResponse(order), nil
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

	return h.dtoMapper.ToCancelOrderResponse(order), nil
}

