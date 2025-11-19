package command

import (
	"context"
	"errors"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.ConfirmOrderCommand = (*confirmOrderCommand)(nil)

type confirmOrderCommand struct {
	orderRepo       orders.Repository
	productRepo     products.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewConfirmOrderCommand(
	orderRepo orders.Repository,
	productRepo products.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.ConfirmOrderCommand {
	return &confirmOrderCommand{
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.confirm_order"),
	}
}

func (c *confirmOrderCommand) Execute(ctx context.Context, orderID string) error {
	// Load order
	orderModel, err := c.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to find order")
		return err
	}

	// Validate stock availability for all items
	for _, item := range orderModel.Items {
		productModel, err := c.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to find product")
			return err
		}

		if productModel.Stock < item.Quantity {
			c.log.Error().
				Str("product_id", item.ProductID).
				Int("requested", item.Quantity).
				Int("available", productModel.Stock).
				Msg("insufficient stock")
			return errors.New("insufficient stock for product: " + item.ProductID)
		}
	}

	// Confirm order (domain logic - adds domain event)
	if err := orderModel.Confirm(); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to confirm order")
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

	// Decrease stock for all products
	for _, item := range orderModel.Items {
		productModel, err := c.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to find product for stock update")
			continue
		}

		if err := productModel.DecreaseStock(item.Quantity); err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to decrease stock")
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
		Msg("order confirmed")

	return nil
}
