package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/pkg/logger"
)

type ProductStockUpdatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewProductStockUpdatedHandler(eventBroker contracts.EventBrokerPublisher) *ProductStockUpdatedHandler {
	return &ProductStockUpdatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("ecommerce.event_handler.product_stock_updated"),
	}
}

func (h *ProductStockUpdatedHandler) Handle(ctx context.Context, domainEvent event_dispatcher.DomainEvent) error {
	// Type assert to ProductStockUpdatedEvent
	stockUpdatedEvent, ok := domainEvent.(product.ProductStockUpdatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in ProductStockUpdatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.ProductStockUpdatedPayload{
		ProductID: stockUpdatedEvent.ProductID,
		OldStock:  stockUpdatedEvent.OldStock,
		NewStock:  stockUpdatedEvent.NewStock,
		UpdatedAt: stockUpdatedEvent.UpdatedAt,
	}

	if err := h.eventBroker.PublishProductStockUpdated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("product_id", stockUpdatedEvent.ProductID).
			Msg("failed to publish ProductStockUpdated event")
		return err
	}

	h.log.Info().
		Str("product_id", stockUpdatedEvent.ProductID).
		Int("old_stock", stockUpdatedEvent.OldStock).
		Int("new_stock", stockUpdatedEvent.NewStock).
		Msg("ProductStockUpdated event published")

	return nil
}

