package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"
)

type OrderCreatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewOrderCreatedHandler(eventBroker contracts.EventBrokerPublisher) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("ecommerce.event_handler.order_created"),
	}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, domainEvent event_dispatcher.DomainEvent) error {
	// Type assert to OrderCreatedEvent
	orderCreatedEvent, ok := domainEvent.(order.OrderCreatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in OrderCreatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.OrderCreatedPayload{
		OrderID:     orderCreatedEvent.OrderID,
		UserID:      orderCreatedEvent.UserID,
		TotalAmount: orderCreatedEvent.TotalAmount,
		ItemCount:   orderCreatedEvent.ItemCount,
		CreatedAt:   orderCreatedEvent.CreatedAt,
	}

	if err := h.eventBroker.PublishOrderCreated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", orderCreatedEvent.OrderID).
			Msg("failed to publish OrderCreated event")
		return err
	}

	h.log.Info().
		Str("order_id", orderCreatedEvent.OrderID).
		Str("user_id", orderCreatedEvent.UserID).
		Msg("OrderCreated event published")

	return nil
}

