-- Rollback: Revert partitioning and remove shard_id field

-- Step 1: Rename tables back
BEGIN;
ALTER TABLE messages RENAME TO messages_partitioned;
ALTER TABLE messages_old RENAME TO messages;
COMMIT;

-- Step 2: Drop partitioned table and all partitions
DROP TABLE IF EXISTS messages_partitioned CASCADE;

-- Step 3: Drop function
DROP FUNCTION IF EXISTS calculate_shard_id(TEXT, TEXT, INT);
