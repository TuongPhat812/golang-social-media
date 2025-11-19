package command

import (
	"context"
	"errors"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.AddOrderItemCommand = (*addOrderItemCommand)(nil)

type addOrderItemCommand struct {
	orderRepo       orders.Repository
	productRepo     products.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewAddOrderItemCommand(
	orderRepo orders.Repository,
	productRepo products.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.AddOrderItemCommand {
	return &addOrderItemCommand{
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.add_order_item"),
	}
}

func (c *addOrderItemCommand) Execute(ctx context.Context, req contracts.AddOrderItemCommandRequest) (order.Order, error) {
	// Load order
	orderModel, err := c.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", req.OrderID).
			Msg("failed to find order")
		return order.Order{}, err
	}

	// Load product to get current price and check availability
	productModel, err := c.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("product_id", req.ProductID).
			Msg("failed to find product")
		return order.Order{}, err
	}

	// Check product availability
	if !productModel.IsAvailable() {
		c.log.Error().
			Str("product_id", req.ProductID).
			Str("status", string(productModel.Status)).
			Int("stock", productModel.Stock).
			Msg("product is not available")
		return order.Order{}, errors.New("product is not available")
	}

	// Check stock availability
	if productModel.Stock < req.Quantity {
		c.log.Error().
			Str("product_id", req.ProductID).
			Int("requested", req.Quantity).
			Int("available", productModel.Stock).
			Msg("insufficient stock")
		return order.Order{}, errors.New("insufficient stock")
	}

	// Create order item
	orderItem, err := order.NewOrderItem(req.ProductID, req.Quantity, productModel.Price)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("product_id", req.ProductID).
			Msg("failed to create order item")
		return order.Order{}, err
	}

	// Add item to order (domain logic)
	if err := orderModel.AddItem(orderItem); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", req.OrderID).
			Msg("failed to add item to order")
		return order.Order{}, err
	}

	// Save domain events BEFORE persisting
	domainEvents := orderModel.Events()

	// Persist updated order
	if err := c.orderRepo.Update(ctx, &orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", req.OrderID).
			Msg("failed to update order")
		return order.Order{}, err
	}

	// Dispatch domain events AFTER successful persistence
	orderModel.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("order_id", req.OrderID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("order_id", req.OrderID).
		Str("product_id", req.ProductID).
		Int("quantity", req.Quantity).
		Msg("order item added")

	return orderModel, nil
}
