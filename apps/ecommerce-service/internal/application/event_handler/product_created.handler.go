package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/pkg/logger"
)

type ProductCreatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewProductCreatedHandler(eventBroker contracts.EventBrokerPublisher) *ProductCreatedHandler {
	return &ProductCreatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("ecommerce.event_handler.product_created"),
	}
}

func (h *ProductCreatedHandler) Handle(ctx context.Context, domainEvent event_dispatcher.DomainEvent) error {
	// Type assert to ProductCreatedEvent
	productCreatedEvent, ok := domainEvent.(product.ProductCreatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in ProductCreatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.ProductCreatedPayload{
		ProductID:   productCreatedEvent.ProductID,
		Name:        productCreatedEvent.Name,
		Description: productCreatedEvent.Description,
		Price:       productCreatedEvent.Price,
		Stock:       productCreatedEvent.Stock,
		CreatedAt:   productCreatedEvent.CreatedAt,
	}

	if err := h.eventBroker.PublishProductCreated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("product_id", productCreatedEvent.ProductID).
			Msg("failed to publish ProductCreated event")
		return err
	}

	h.log.Info().
		Str("product_id", productCreatedEvent.ProductID).
		Str("name", productCreatedEvent.Name).
		Msg("ProductCreated event published")

	return nil
}

