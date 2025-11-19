package command

import (
	"context"
	"errors"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/unit_of_work"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.ConfirmOrderCommand = (*confirmOrderCommand)(nil)

type confirmOrderCommand struct {
	uowFactory      unit_of_work.Factory
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewConfirmOrderCommand(
	uowFactory unit_of_work.Factory,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.ConfirmOrderCommand {
	return &confirmOrderCommand{
		uowFactory:      uowFactory,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.confirm_order"),
	}
}

func (c *confirmOrderCommand) Execute(ctx context.Context, orderID string) error {
	// Create unit of work
	uow, err := c.uowFactory.New(ctx)
	if err != nil {
		c.log.Error().
			Err(err).
			Msg("failed to create unit of work")
		return err
	}
	defer uow.Rollback() // Ensure rollback if commit fails

	// Load order using UoW repository
	orderModel, err := uow.Orders().FindByID(ctx, orderID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to find order")
		return err
	}

	// Validate stock availability for all items
	for _, item := range orderModel.Items {
		productModel, err := uow.Products().FindByID(ctx, item.ProductID)
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

	// Update order status using UoW repository
	if err := uow.Orders().Update(ctx, &orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to update order")
		return err
	}

	// Decrease stock for all products using UoW repository
	for _, item := range orderModel.Items {
		productModel, err := uow.Products().FindByID(ctx, item.ProductID)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to find product for stock update")
			return err
		}

		if err := productModel.DecreaseStock(item.Quantity); err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to decrease stock")
			return err
		}

		// Save product stock events
		productEvents := productModel.Events()

		if err := uow.Products().Update(ctx, &productModel); err != nil {
			c.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Msg("failed to update product stock")
			return err
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

	// Commit transaction
	if err := uow.Commit(); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to commit transaction")
		return err
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
