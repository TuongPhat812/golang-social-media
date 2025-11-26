package outbox

import (
	"context"
	"encoding/json"

	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
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
		log:        logger.Component("auth.outbox"),
	}
}

// SaveEvent saves a domain event to the outbox
func (s *OutboxService) SaveEvent(ctx context.Context, event interface{}) error {
	// Extract event information
	var aggregateID, aggregateType, eventType string
	var eventVersion int = 1 // Default version

	// Try to get event type from interface
	if domainEvent, ok := event.(interface {
		Type() string
	}); ok {
		eventType = domainEvent.Type()
	}

	// Extract aggregate info from struct fields
	eventMap, err := structToMap(event)
	if err != nil {
		return err
	}

	// Try common field names for auth-service events
	if id, ok := eventMap["UserID"].(string); ok && id != "" {
		aggregateID = id
		aggregateType = "User"
	} else if id, ok := eventMap["RoleID"].(string); ok && id != "" {
		aggregateID = id
		aggregateType = "Role"
	} else if id, ok := eventMap["PermissionID"].(string); ok && id != "" {
		aggregateID = id
		aggregateType = "Permission"
	}

	// If event type is still empty, try to infer from struct name
	if eventType == "" {
		eventType = inferEventType(event)
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

// inferEventType tries to infer event type from struct name or fields
func inferEventType(event interface{}) string {
	eventMap, err := structToMap(event)
	if err != nil {
		return "Unknown"
	}
	// Try to get from common patterns
	if _, ok := eventMap["UserID"]; ok {
		return "UserEvent"
	}
	if _, ok := eventMap["RoleID"]; ok {
		return "RoleEvent"
	}
	if _, ok := eventMap["PermissionID"]; ok {
		return "PermissionEvent"
	}
	return "Unknown"
}

