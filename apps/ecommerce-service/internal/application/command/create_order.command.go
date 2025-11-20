package command

import (
	"context"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	"golang-social-media/apps/ecommerce-service/internal/application/unit_of_work"
	"golang-social-media/apps/ecommerce-service/internal/domain/factories"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/eventstore"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/outbox"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.CreateOrderCommand = (*createOrderCommand)(nil)

type createOrderCommand struct {
	uowFactory    unit_of_work.Factory
	orderFactory  factories.OrderFactory
	outboxService *outbox.OutboxService
	eventStore    *eventstore.EventStoreService
	log           *zerolog.Logger
}

func NewCreateOrderCommand(
	uowFactory unit_of_work.Factory,
	orderFactory factories.OrderFactory,
	outboxService *outbox.OutboxService,
	eventStore *eventstore.EventStoreService,
) contracts.CreateOrderCommand {
	return &createOrderCommand{
		uowFactory:    uowFactory,
		orderFactory:  orderFactory,
		outboxService: outboxService,
		eventStore:    eventStore,
		log:           logger.Component("ecommerce.command.create_order"),
	}
}

func (c *createOrderCommand) Execute(ctx context.Context, req contracts.CreateOrderCommandRequest) (order.Order, error) {
	// Use factory to create order
	orderModel, err := c.orderFactory.CreateOrder(req.UserID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to create order using factory")
		return order.Order{}, err
	}

	// Create unit of work
	uow, err := c.uowFactory.New(ctx)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to create unit of work")
		return order.Order{}, err
	}
	defer uow.Rollback() // Ensure rollback if commit fails

	// Save domain events BEFORE persisting
	domainEvents := orderModel.Events()

	// Persist to database using UoW repository
	if err := uow.Orders().Create(ctx, orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to persist order")
		return order.Order{}, err
	}

	// Save events to outbox and event store within the same transaction
	for _, domainEvent := range domainEvents {
		// Save to outbox for reliable publishing
		if err := c.outboxService.SaveEvent(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("order_id", orderModel.ID).
				Msg("failed to save event to outbox")
			return order.Order{}, err
		}

		// Save to event store for event sourcing
		metadata := map[string]interface{}{
			"user_id": orderModel.UserID,
		}
		if err := c.eventStore.Append(ctx, domainEvent, metadata); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("order_id", orderModel.ID).
				Msg("failed to append event to event store")
			return order.Order{}, err
		}
	}

	// Commit transaction (includes outbox and event store writes)
	if err := uow.Commit(); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", orderModel.ID).
			Msg("failed to commit transaction")
		return order.Order{}, err
	}

	// Clear events after successful persistence
	orderModel.ClearEvents()

	c.log.Info().
		Str("order_id", orderModel.ID).
		Str("user_id", orderModel.UserID).
		Int("event_count", len(domainEvents)).
		Msg("order created with events saved to outbox and event store")

	return *orderModel, nil
}

