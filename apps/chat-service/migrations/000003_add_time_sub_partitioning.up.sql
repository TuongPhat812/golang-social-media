-- Migration: Add time-based sub-partitioning to existing hash partitions
-- This creates a 2-level partitioning: Hash (64) × Range (monthly)
-- Only run this if you need time-based data lifecycle management
-- 
-- WARNING: This migration is complex and requires careful planning
-- Consider if you really need sub-partitioning before running this

-- Step 1: Create helper function to get partition name
CREATE OR REPLACE FUNCTION get_time_partition_name(base_name TEXT, year INT, month INT)
RETURNS TEXT AS $$
BEGIN
    RETURN format('%s_%s_%02d', base_name, year, month);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Step 2: Convert each hash partition to sub-partitioned
-- This example converts messages_p0 to have monthly sub-partitions
-- You would need to do this for all 64 partitions (p0-p63)

-- Example for messages_p0:
-- First, we need to recreate messages_p0 as a sub-partitioned table
-- This is complex and requires:
-- 1. Create new sub-partitioned table
-- 2. Migrate data
-- 3. Drop old partition
-- 4. Attach new partition

-- NOTE: This is a template. Actual implementation requires:
-- - Careful data migration
-- - Downtime or careful online migration
-- - Testing on staging first

-- Example structure (DO NOT RUN DIRECTLY - needs careful implementation):
/*
-- For each hash partition (p0-p63), convert to sub-partitioned:
DO $$
DECLARE
    partition_num INT;
    year_val INT;
    month_val INT;
    partition_name TEXT;
    sub_partition_name TEXT;
BEGIN
    FOR partition_num IN 0..63 LOOP
        partition_name := format('messages_p%s', partition_num);
        
        -- Create monthly partitions for current and next 12 months
        FOR year_val IN 2024..2025 LOOP
            FOR month_val IN 1..12 LOOP
                sub_partition_name := format('%s_%s_%02d', partition_name, year_val, month_val);
                
                -- Create sub-partition
                EXECUTE format(
                    'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
                    sub_partition_name,
                    partition_name,
                    format('%s-%02d-01', year_val, month_val),
                    CASE 
                        WHEN month_val = 12 THEN format('%s-01-01', year_val + 1)
                        ELSE format('%s-%02d-01', year_val, month_val + 1)
                    END
                );
            END LOOP;
        END LOOP;
    END LOOP;
END $$;
*/

-- IMPORTANT: This migration is provided as a reference
-- Actual implementation should be done carefully with:
-- 1. Backup first
-- 2. Test on staging
-- 3. Plan for downtime or use online migration tools
-- 4. Monitor performance impact

