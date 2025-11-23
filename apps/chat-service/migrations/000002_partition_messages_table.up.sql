-- Migration: Partition messages table by hash of (sender_id, receiver_id)
-- This reduces contention by distributing writes across multiple partitions
-- Each partition handles ~1/64 of the total load
-- Includes shard_id field for easy migration when changing partition count
-- Uses 64 partitions for maximum write distribution

-- Step 1: Create function to calculate shard_id from (sender_id, receiver_id)
-- This function will be used to calculate which partition a message belongs to
CREATE OR REPLACE FUNCTION calculate_shard_id(sender_id TEXT, receiver_id TEXT, num_shards INT DEFAULT 64)
RETURNS INT AS $$
BEGIN
    -- Use consistent hash of (sender_id, receiver_id) pair
    -- Normalize: always use smaller ID first for consistent hashing
    DECLARE
        pair TEXT;
        hash_val BIGINT;
    BEGIN
        -- Normalize pair: always use lexicographically smaller ID first
        IF sender_id < receiver_id THEN
            pair := sender_id || ':' || receiver_id;
        ELSE
            pair := receiver_id || ':' || sender_id;
        END IF;
        
        -- Calculate hash using PostgreSQL's hash function
        hash_val := abs(hashtext(pair));
        
        -- Return modulo
        RETURN hash_val % num_shards;
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Step 2: Create new partitioned table structure with shard_id field
CREATE TABLE IF NOT EXISTS messages_new (
    id UUID NOT NULL,
    sender_id TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    shard_id INT NOT NULL GENERATED ALWAYS AS (
        calculate_shard_id(sender_id, receiver_id, 64)
    ) STORED,
    PRIMARY KEY (id, sender_id, receiver_id)
) PARTITION BY HASH (sender_id, receiver_id);

-- Step 3: Create 64 partitions for maximum write distribution
-- Each partition will handle approximately 1/64 of the writes
CREATE TABLE IF NOT EXISTS messages_p0 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 0);
CREATE TABLE IF NOT EXISTS messages_p1 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 1);
CREATE TABLE IF NOT EXISTS messages_p2 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 2);
CREATE TABLE IF NOT EXISTS messages_p3 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 3);
CREATE TABLE IF NOT EXISTS messages_p4 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 4);
CREATE TABLE IF NOT EXISTS messages_p5 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 5);
CREATE TABLE IF NOT EXISTS messages_p6 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 6);
CREATE TABLE IF NOT EXISTS messages_p7 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 7);
CREATE TABLE IF NOT EXISTS messages_p8 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 8);
CREATE TABLE IF NOT EXISTS messages_p9 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 9);
CREATE TABLE IF NOT EXISTS messages_p10 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 10);
CREATE TABLE IF NOT EXISTS messages_p11 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 11);
CREATE TABLE IF NOT EXISTS messages_p12 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 12);
CREATE TABLE IF NOT EXISTS messages_p13 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 13);
CREATE TABLE IF NOT EXISTS messages_p14 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 14);
CREATE TABLE IF NOT EXISTS messages_p15 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 15);
CREATE TABLE IF NOT EXISTS messages_p16 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 16);
CREATE TABLE IF NOT EXISTS messages_p17 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 17);
CREATE TABLE IF NOT EXISTS messages_p18 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 18);
CREATE TABLE IF NOT EXISTS messages_p19 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 19);
CREATE TABLE IF NOT EXISTS messages_p20 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 20);
CREATE TABLE IF NOT EXISTS messages_p21 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 21);
CREATE TABLE IF NOT EXISTS messages_p22 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 22);
CREATE TABLE IF NOT EXISTS messages_p23 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 23);
CREATE TABLE IF NOT EXISTS messages_p24 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 24);
CREATE TABLE IF NOT EXISTS messages_p25 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 25);
CREATE TABLE IF NOT EXISTS messages_p26 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 26);
CREATE TABLE IF NOT EXISTS messages_p27 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 27);
CREATE TABLE IF NOT EXISTS messages_p28 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 28);
CREATE TABLE IF NOT EXISTS messages_p29 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 29);
CREATE TABLE IF NOT EXISTS messages_p30 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 30);
CREATE TABLE IF NOT EXISTS messages_p31 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 31);
CREATE TABLE IF NOT EXISTS messages_p32 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 32);
CREATE TABLE IF NOT EXISTS messages_p33 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 33);
CREATE TABLE IF NOT EXISTS messages_p34 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 34);
CREATE TABLE IF NOT EXISTS messages_p35 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 35);
CREATE TABLE IF NOT EXISTS messages_p36 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 36);
CREATE TABLE IF NOT EXISTS messages_p37 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 37);
CREATE TABLE IF NOT EXISTS messages_p38 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 38);
CREATE TABLE IF NOT EXISTS messages_p39 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 39);
CREATE TABLE IF NOT EXISTS messages_p40 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 40);
CREATE TABLE IF NOT EXISTS messages_p41 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 41);
CREATE TABLE IF NOT EXISTS messages_p42 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 42);
CREATE TABLE IF NOT EXISTS messages_p43 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 43);
CREATE TABLE IF NOT EXISTS messages_p44 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 44);
CREATE TABLE IF NOT EXISTS messages_p45 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 45);
CREATE TABLE IF NOT EXISTS messages_p46 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 46);
CREATE TABLE IF NOT EXISTS messages_p47 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 47);
CREATE TABLE IF NOT EXISTS messages_p48 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 48);
CREATE TABLE IF NOT EXISTS messages_p49 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 49);
CREATE TABLE IF NOT EXISTS messages_p50 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 50);
CREATE TABLE IF NOT EXISTS messages_p51 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 51);
CREATE TABLE IF NOT EXISTS messages_p52 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 52);
CREATE TABLE IF NOT EXISTS messages_p53 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 53);
CREATE TABLE IF NOT EXISTS messages_p54 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 54);
CREATE TABLE IF NOT EXISTS messages_p55 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 55);
CREATE TABLE IF NOT EXISTS messages_p56 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 56);
CREATE TABLE IF NOT EXISTS messages_p57 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 57);
CREATE TABLE IF NOT EXISTS messages_p58 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 58);
CREATE TABLE IF NOT EXISTS messages_p59 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 59);
CREATE TABLE IF NOT EXISTS messages_p60 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 60);
CREATE TABLE IF NOT EXISTS messages_p61 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 61);
CREATE TABLE IF NOT EXISTS messages_p62 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 62);
CREATE TABLE IF NOT EXISTS messages_p63 PARTITION OF messages_new FOR VALUES WITH (MODULUS 64, REMAINDER 63);

-- Step 4: Create indexes on each partition for better query performance
-- We'll create indexes on all 64 partitions for sender_id, receiver_id, created_at, and shard_id
-- Using a loop-like approach with DO block to generate indexes dynamically
DO $$
DECLARE
    i INT;
BEGIN
    FOR i IN 0..63 LOOP
        -- Index on sender_id
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_messages_p%s_sender_id ON messages_p%s(sender_id)', i, i);
        -- Index on receiver_id
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_messages_p%s_receiver_id ON messages_p%s(receiver_id)', i, i);
        -- Index on created_at for time-based queries
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_messages_p%s_created_at ON messages_p%s(created_at DESC)', i, i);
        -- Index on shard_id for direct partition queries
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_messages_p%s_shard_id ON messages_p%s(shard_id)', i, i);
    END LOOP;
END $$;

-- Step 5: Migrate existing data (if any)
-- Calculate shard_id for existing data and insert into new partitioned table
INSERT INTO messages_new (id, sender_id, receiver_id, content, created_at)
SELECT 
    id, 
    sender_id, 
    receiver_id, 
    content, 
    created_at
FROM messages
ON CONFLICT DO NOTHING;

-- Step 6: Rename tables (atomic operation)
BEGIN;
ALTER TABLE messages RENAME TO messages_old;
ALTER TABLE messages_new RENAME TO messages;
COMMIT;

-- Step 7: Drop old table (optional - can keep for backup)
-- DROP TABLE IF EXISTS messages_old;
