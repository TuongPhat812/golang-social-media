package event_dispatcher

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/pkg/logger"
)

// EventHandler handles a specific domain event type
type EventHandler interface {
	Handle(ctx context.Context, event notification.DomainEvent) error
}

// Dispatcher dispatches domain events to their respective handlers
type Dispatcher struct {
	handlers map[string][]EventHandler
	log      *zerolog.Logger
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
		log:      logger.Component("notification.event_dispatcher"),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (d *Dispatcher) RegisterHandler(eventType string, handler EventHandler) {
	if d.handlers[eventType] == nil {
		d.handlers[eventType] = make([]EventHandler, 0)
	}
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch dispatches a domain event to all registered handlers
func (d *Dispatcher) Dispatch(ctx context.Context, event notification.DomainEvent) error {
	eventType := event.Type()
	handlers := d.handlers[eventType]

	if len(handlers) == 0 {
		d.log.Warn().
			Str("event_type", eventType).
			Msg("no handlers registered for event type")
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			d.log.Error().
				Err(err).
				Str("event_type", eventType).
				Msg("failed to handle domain event")
			// Continue with other handlers even if one fails
			// In production, you might want to use an outbox pattern for retry
		}
	}

	return nil
}

