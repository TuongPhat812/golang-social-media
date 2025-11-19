package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"
)

type OrderCancelledHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewOrderCancelledHandler(eventBroker contracts.EventBrokerPublisher) *OrderCancelledHandler {
	return &OrderCancelledHandler{
		eventBroker: eventBroker,
		log:         logger.Component("ecommerce.event_handler.order_cancelled"),
	}
}

func (h *OrderCancelledHandler) Handle(ctx context.Context, domainEvent event_dispatcher.DomainEvent) error {
	// Type assert to OrderCancelledEvent
	orderCancelledEvent, ok := domainEvent.(order.OrderCancelledEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in OrderCancelledHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.OrderCancelledPayload{
		OrderID:    orderCancelledEvent.OrderID,
		UserID:     orderCancelledEvent.UserID,
		CancelledAt: orderCancelledEvent.CancelledAt,
	}

	if err := h.eventBroker.PublishOrderCancelled(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", orderCancelledEvent.OrderID).
			Msg("failed to publish OrderCancelled event")
		return err
	}

	h.log.Info().
		Str("order_id", orderCancelledEvent.OrderID).
		Str("user_id", orderCancelledEvent.UserID).
		Msg("OrderCancelled event published")

	return nil
}

