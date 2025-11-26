package postgres

import (
	"time"
)

// EventStoreModel represents an event in the event store
type EventStoreModel struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey"`
	AggregateID   string    `gorm:"column:aggregate_id;type:text;not null"`
	AggregateType string    `gorm:"column:aggregate_type;type:text;not null"`
	EventType     string    `gorm:"column:event_type;type:text;not null"`
	EventVersion  int       `gorm:"column:event_version;type:integer;not null;default:1"`
	Payload       string    `gorm:"column:payload;type:jsonb;not null"` // JSON string
	Metadata      *string   `gorm:"column:metadata;type:jsonb"`         // JSON string, optional
	OccurredAt    time.Time `gorm:"column:occurred_at;not null"`
}

func (EventStoreModel) TableName() string {
	return "event_store"
}

