package event_dispatcher

import (
	"context"

	"golang-social-media/apps/auth-service/internal/domain/user"
)

// EventHandler handles domain events
type EventHandler interface {
	Handle(ctx context.Context, domainEvent user.DomainEvent) error
}

// Dispatcher dispatches domain events to registered handlers
type Dispatcher struct {
	handlers map[string][]EventHandler
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (d *Dispatcher) RegisterHandler(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch dispatches a domain event to all registered handlers
func (d *Dispatcher) Dispatch(ctx context.Context, domainEvent user.DomainEvent) error {
	eventType := domainEvent.Type()
	handlers := d.handlers[eventType]

	for _, handler := range handlers {
		if err := handler.Handle(ctx, domainEvent); err != nil {
			return err
		}
	}

	return nil
}

