package command

import (
	"context"
	"time"

	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var _ contracts.CreateOrderCommand = (*createOrderCommand)(nil)

type createOrderCommand struct {
	repo            orders.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewCreateOrderCommand(
	repo orders.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CreateOrderCommand {
	return &createOrderCommand{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("ecommerce.command.create_order"),
	}
}

func (c *createOrderCommand) Execute(ctx context.Context, req contracts.CreateOrderCommandRequest) (order.Order, error) {
	createdAt := time.Now().UTC()
	orderModel := order.Order{
		ID:          uuid.NewString(),
		UserID:      req.UserID,
		Status:      order.StatusDraft,
		Items:       []order.OrderItem{},
		TotalAmount: 0,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	// Validate business rules
	if err := orderModel.Validate(); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("order validation failed")
		return order.Order{}, err
	}

	// Domain logic: create order (this adds domain events internally)
	orderModel.Create()

	// Save domain events BEFORE persisting
	domainEvents := orderModel.Events()

	// Persist to database
	if err := c.repo.Create(ctx, &orderModel); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to persist order")
		return order.Order{}, err
	}

	// Dispatch domain events AFTER successful persistence
	orderModel.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("order_id", orderModel.ID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("order_id", orderModel.ID).
		Str("user_id", orderModel.UserID).
		Msg("order created")

	return orderModel, nil
}
