-- Create event store table for event sourcing
CREATE TABLE IF NOT EXISTS event_store (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_version INTEGER NOT NULL DEFAULT 1,
    payload JSONB NOT NULL,
    metadata JSONB,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    INDEX idx_event_store_aggregate (aggregate_id, aggregate_type),
    INDEX idx_event_store_occurred_at (occurred_at)
);

