package event_dispatcher

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/pkg/logger"
)

// EventHandler handles domain events
type EventHandler interface {
	Handle(ctx context.Context, domainEvent message.DomainEvent) error
}

// Dispatcher dispatches domain events to registered handlers
type Dispatcher struct {
	handlers map[string][]EventHandler
	log      *zerolog.Logger
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
		log:      logger.Component("chat.event_dispatcher"),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (d *Dispatcher) RegisterHandler(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch dispatches a domain event to all registered handlers
func (d *Dispatcher) Dispatch(ctx context.Context, domainEvent message.DomainEvent) error {
	eventType := domainEvent.Type()
	handlers := d.handlers[eventType]

	if len(handlers) == 0 {
		d.log.Warn().
			Str("event_type", eventType).
			Msg("no handlers registered for event type")
		return nil
	}

	d.log.Info().
		Str("event_type", eventType).
		Int("handler_count", len(handlers)).
		Msg("dispatching domain event")

	for _, handler := range handlers {
		if err := handler.Handle(ctx, domainEvent); err != nil {
			d.log.Error().
				Err(err).
				Str("event_type", eventType).
				Msg("handler failed to process event")
			return err
		}
	}

	d.log.Info().
		Str("event_type", eventType).
		Msg("domain event dispatched successfully")

	return nil
}

