-- Drop event store table
DROP INDEX IF EXISTS idx_event_store_occurred_at;
DROP INDEX IF EXISTS idx_event_store_aggregate;
DROP TABLE IF EXISTS event_store;

