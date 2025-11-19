package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"
)

type OrderConfirmedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewOrderConfirmedHandler(eventBroker contracts.EventBrokerPublisher) *OrderConfirmedHandler {
	return &OrderConfirmedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("ecommerce.event_handler.order_confirmed"),
	}
}

func (h *OrderConfirmedHandler) Handle(ctx context.Context, domainEvent event_dispatcher.DomainEvent) error {
	// Type assert to OrderConfirmedEvent
	orderConfirmedEvent, ok := domainEvent.(order.OrderConfirmedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in OrderConfirmedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.OrderConfirmedPayload{
		OrderID:     orderConfirmedEvent.OrderID,
		UserID:      orderConfirmedEvent.UserID,
		TotalAmount: orderConfirmedEvent.TotalAmount,
		ItemCount:   orderConfirmedEvent.ItemCount,
		ConfirmedAt: orderConfirmedEvent.ConfirmedAt,
	}

	if err := h.eventBroker.PublishOrderConfirmed(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("order_id", orderConfirmedEvent.OrderID).
			Msg("failed to publish OrderConfirmed event")
		return err
	}

	h.log.Info().
		Str("order_id", orderConfirmedEvent.OrderID).
		Str("user_id", orderConfirmedEvent.UserID).
		Msg("OrderConfirmed event published")

	return nil
}

