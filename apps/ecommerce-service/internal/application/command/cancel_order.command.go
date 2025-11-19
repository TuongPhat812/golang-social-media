package command

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/unit_of_work"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.CancelOrderCommand = (*cancelOrderCommand)(nil)

type cancelOrderCommand struct {
	uowFactory      unit_of_work.Factory
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewCancelOrderCommand(
	uowFactory unit_of_work.Factory,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CancelOrderCommand {
	return &cancelOrderCommand{
		uowFactory:      uowFactory,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.cancel_order"),
	}
}

func (c *cancelOrderCommand) Execute(ctx context.Context, orderID string) error {
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

	// Update order status using UoW repository
	if err := uow.Orders().Update(ctx, &orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to update order")
		return err
	}

	// If order was confirmed, restore stock using UoW repository
	if orderModel.Status == order.StatusCancelled {
		for _, item := range orderModel.Items {
			productModel, err := uow.Products().FindByID(ctx, item.ProductID)
			if err != nil {
				c.log.Error().
					Err(err).
					Str("product_id", item.ProductID).
					Msg("failed to find product for stock restoration")
				return err
			}

			if err := productModel.IncreaseStock(item.Quantity); err != nil {
				c.log.Error().
					Err(err).
					Str("product_id", item.ProductID).
					Msg("failed to increase stock")
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
		Msg("order cancelled")

	return nil
}
