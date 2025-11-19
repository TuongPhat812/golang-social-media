package outbox

import (
	"context"
	"encoding/json"

	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

// OutboxService handles outbox operations
type OutboxService struct {
	outboxRepo *postgres.OutboxRepository
	log        *zerolog.Logger
}

// NewOutboxService creates a new OutboxService
func NewOutboxService(outboxRepo *postgres.OutboxRepository) *OutboxService {
	return &OutboxService{
		outboxRepo: outboxRepo,
		log:        logger.Component("ecommerce.outbox"),
	}
}

// SaveEvent saves a domain event to the outbox
func (s *OutboxService) SaveEvent(ctx context.Context, event interface{}) error {
	// Extract event information using type assertion
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

	// Save to outbox
	err := s.outboxRepo.Create(ctx, aggregateID, aggregateType, eventType, eventVersion, event)
	if err != nil {
		s.log.Error().
			Err(err).
			Str("event_type", eventType).
			Str("aggregate_id", aggregateID).
			Msg("failed to save event to outbox")
		return err
	}

	s.log.Debug().
		Str("event_type", eventType).
		Str("aggregate_id", aggregateID).
		Int("version", eventVersion).
		Msg("event saved to outbox")

	return nil
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

