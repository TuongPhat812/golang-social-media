package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventStoreRepository handles event store operations
type EventStoreRepository struct {
	db *gorm.DB
}

// NewEventStoreRepository creates a new EventStoreRepository
func NewEventStoreRepository(db *gorm.DB) *EventStoreRepository {
	return &EventStoreRepository{db: db}
}

// Append stores an event in the event store
func (r *EventStoreRepository) Append(ctx context.Context, aggregateID, aggregateType, eventType string, eventVersion int, payload interface{}, metadata map[string]interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var metadataJSON *string
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		metadataStr := string(metadataBytes)
		metadataJSON = &metadataStr
	}

	model := EventStoreModel{
		ID:            uuid.NewString(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		EventVersion:  eventVersion,
		Payload:       string(payloadJSON),
		Metadata:      metadataJSON,
		OccurredAt:    time.Now().UTC(),
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

// GetEventsByAggregate retrieves all events for a specific aggregate
func (r *EventStoreRepository) GetEventsByAggregate(ctx context.Context, aggregateID, aggregateType string) ([]EventStoreModel, error) {
	var events []EventStoreModel
	err := r.db.WithContext(ctx).
		Where("aggregate_id = ? AND aggregate_type = ?", aggregateID, aggregateType).
		Order("occurred_at ASC").
		Find(&events).Error
	return events, err
}

// GetEventsByType retrieves all events of a specific type
func (r *EventStoreRepository) GetEventsByType(ctx context.Context, eventType string, limit int) ([]EventStoreModel, error) {
	var events []EventStoreModel
	err := r.db.WithContext(ctx).
		Where("event_type = ?", eventType).
		Order("occurred_at DESC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// GetEventsByTimeRange retrieves events within a time range
func (r *EventStoreRepository) GetEventsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit int) ([]EventStoreModel, error) {
	var events []EventStoreModel
	err := r.db.WithContext(ctx).
		Where("occurred_at >= ? AND occurred_at <= ?", startTime, endTime).
		Order("occurred_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

