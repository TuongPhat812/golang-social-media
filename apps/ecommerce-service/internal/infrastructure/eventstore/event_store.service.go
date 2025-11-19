package eventstore

import (
	"context"
	"encoding/json"
	"time"

	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

// EventStoreService handles event store operations
type EventStoreService struct {
	eventStoreRepo *postgres.EventStoreRepository
	log            *zerolog.Logger
}

// NewEventStoreService creates a new EventStoreService
func NewEventStoreService(eventStoreRepo *postgres.EventStoreRepository) *EventStoreService {
	return &EventStoreService{
		eventStoreRepo: eventStoreRepo,
		log:            logger.Component("ecommerce.event_store"),
	}
}

// Append stores an event in the event store
func (s *EventStoreService) Append(ctx context.Context, event interface{}, metadata map[string]interface{}) error {
	// Extract event information
	var aggregateID, aggregateType, eventType string
	var eventVersion int

	// Try to get event metadata from common interface
	if domainEvent, ok := event.(interface {
		Type() string
		Version() int
		AggregateID() string
		AggregateType() string
	}); ok {
		eventType = domainEvent.Type()
		eventVersion = domainEvent.Version()
		aggregateID = domainEvent.AggregateID()
		aggregateType = domainEvent.AggregateType()
	} else {
		// Fallback: try to extract from struct fields
		eventMap, err := structToMap(event)
		if err != nil {
			return err
		}
		// Try common field names
		if id, ok := eventMap["ProductID"].(string); ok {
			aggregateID = id
			aggregateType = "Product"
		} else if id, ok := eventMap["OrderID"].(string); ok {
			aggregateID = id
			aggregateType = "Order"
		}
		eventVersion = 1 // Default version
	}

	// Store in event store
	err := s.eventStoreRepo.Append(ctx, aggregateID, aggregateType, eventType, eventVersion, event, metadata)
	if err != nil {
		s.log.Error().
			Err(err).
			Str("event_type", eventType).
			Str("aggregate_id", aggregateID).
			Msg("failed to append event to event store")
		return err
	}

	s.log.Debug().
		Str("event_type", eventType).
		Str("aggregate_id", aggregateID).
		Int("version", eventVersion).
		Msg("event appended to event store")

	return nil
}

// GetEventsByAggregate retrieves all events for a specific aggregate
func (s *EventStoreService) GetEventsByAggregate(ctx context.Context, aggregateID, aggregateType string) ([]postgres.EventStoreModel, error) {
	return s.eventStoreRepo.GetEventsByAggregate(ctx, aggregateID, aggregateType)
}

// GetEventsByType retrieves all events of a specific type
func (s *EventStoreService) GetEventsByType(ctx context.Context, eventType string, limit int) ([]postgres.EventStoreModel, error) {
	return s.eventStoreRepo.GetEventsByType(ctx, eventType, limit)
}

// GetEventsByTimeRange retrieves events within a time range
func (s *EventStoreService) GetEventsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit int) ([]postgres.EventStoreModel, error) {
	return s.eventStoreRepo.GetEventsByTimeRange(ctx, startTime, endTime, limit)
}

// structToMap converts a struct to a map
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

