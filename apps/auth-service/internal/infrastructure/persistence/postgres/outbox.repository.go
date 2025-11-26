package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OutboxRepository handles outbox operations
type OutboxRepository struct {
	db *gorm.DB
}

// NewOutboxRepository creates a new OutboxRepository
func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

// NewOutboxRepositoryWithTx creates an OutboxRepository with a specific transaction
func NewOutboxRepositoryWithTx(tx *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: tx}
}

// Create stores an event in the outbox
func (r *OutboxRepository) Create(ctx context.Context, aggregateID, aggregateType, eventType string, eventVersion int, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	model := OutboxModel{
		ID:            uuid.NewString(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		EventVersion:  eventVersion,
		Payload:       string(payloadJSON),
		Status:        OutboxStatusPending,
		RetryCount:    0,
		CreatedAt:     time.Now().UTC(),
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

// GetPendingEvents retrieves pending events from the outbox
func (r *OutboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]OutboxModel, error) {
	var events []OutboxModel
	err := r.db.WithContext(ctx).
		Where("status = ?", OutboxStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// MarkAsPublished marks an event as published
func (r *OutboxRepository) MarkAsPublished(ctx context.Context, id string) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&OutboxModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       OutboxStatusPublished,
			"published_at": now,
		}).Error
}

// MarkAsFailed marks an event as failed and increments retry count
func (r *OutboxRepository) MarkAsFailed(ctx context.Context, id string, errorMessage string) error {
	return r.db.WithContext(ctx).
		Model(&OutboxModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        OutboxStatusFailed,
			"retry_count":   gorm.Expr("retry_count + 1"),
			"error_message": errorMessage,
		}).Error
}

// IncrementRetry increments the retry count for an event
func (r *OutboxRepository) IncrementRetry(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&OutboxModel{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + 1")).Error
}

// ResetToPending resets an event back to pending status for retry
func (r *OutboxRepository) ResetToPending(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&OutboxModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       OutboxStatusPending,
			"error_message": nil,
		}).Error
}

