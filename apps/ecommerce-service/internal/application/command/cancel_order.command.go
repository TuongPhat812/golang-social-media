package command

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.CancelOrderCommand = (*cancelOrderCommand)(nil)

type cancelOrderCommand struct {
	orderRepo       orders.Repository
	productRepo     products.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewCancelOrderCommand(
	orderRepo orders.Repository,
	productRepo products.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CancelOrderCommand {
	return &cancelOrderCommand{
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.cancel_order"),
	}
}

func (c *cancelOrderCommand) Execute(ctx context.Context, orderID string) error {
	// Load order
	orderModel, err := c.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to find order")
		return err
	}

	// Cancel order (domain logic - adds domain event)
	if err := orderModel.Cancel(); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to cancel order")
		return err
	}

	// Save domain events BEFORE persisting
	domainEvents := orderModel.Events()

	// Update order status
	if err := c.orderRepo.Update(ctx, &orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to update order")
		return err
	}

	// If order was confirmed, restore stock
	if orderModel.Status == order.StatusCancelled {
		for _, item := range orderModel.Items {
			productModel, err := c.productRepo.FindByID(ctx, item.ProductID)
			if err != nil {
				c.log.Error().
					Err(err).
					Str("product_id", item.ProductID).
					Msg("failed to find product for stock restoration")
				continue
			}

			if err := productModel.IncreaseStock(item.Quantity); err != nil {
				c.log.Error().
					Err(err).
					Str("product_id", item.ProductID).
					Msg("failed to increase stock")
				continue
			}

			// Save product stock events
			productEvents := productModel.Events()

			if err := c.productRepo.Update(ctx, &productModel); err != nil {
				c.log.Error().
					Err(err).
					Str("product_id", item.ProductID).
					Msg("failed to update product stock")
				continue
			}

			// Dispatch product stock events
			productModel.ClearEvents()
			for _, domainEvent := range productEvents {
				if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
					c.log.Error().
						Err(err).
						Str("event_type", domainEvent.Type()).
						Str("product_id", item.ProductID).
						Msg("failed to dispatch product stock event")
				}
			}
		}
	}

	// Dispatch order events AFTER successful persistence
	orderModel.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("order_id", orderID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("order_id", orderID).
		Str("user_id", orderModel.UserID).
		Msg("order cancelled")

	return nil
}
