package postgres

import (
	"time"
)

// OutboxStatus represents the status of an outbox event
type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusPublished OutboxStatus = "published"
	OutboxStatusFailed    OutboxStatus = "failed"
)

// OutboxModel represents an outbox event in the database
type OutboxModel struct {
	ID            string       `gorm:"column:id;type:uuid;primaryKey"`
	AggregateID   string       `gorm:"column:aggregate_id;type:text;not null"`
	AggregateType string       `gorm:"column:aggregate_type;type:text;not null"`
	EventType     string       `gorm:"column:event_type;type:text;not null"`
	EventVersion  int          `gorm:"column:event_version;type:integer;not null;default:1"`
	Payload       string       `gorm:"column:payload;type:jsonb;not null"` // JSON string
	Status        OutboxStatus `gorm:"column:status;type:text;not null;default:'pending'"`
	RetryCount    int          `gorm:"column:retry_count;type:integer;not null;default:0"`
	CreatedAt     time.Time    `gorm:"column:created_at;not null"`
	PublishedAt   *time.Time   `gorm:"column:published_at"`
	ErrorMessage  *string      `gorm:"column:error_message;type:text"`
}

func (OutboxModel) TableName() string {
	return "outbox"
}

