-- Drop outbox table
DROP INDEX IF EXISTS idx_outbox_created_at;
DROP INDEX IF EXISTS idx_outbox_status;
DROP TABLE IF EXISTS outbox;

