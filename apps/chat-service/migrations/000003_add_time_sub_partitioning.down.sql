-- Rollback: Remove time-based sub-partitioning
-- This reverts sub-partitions back to simple hash partitions
-- 
-- WARNING: This is complex and may require data migration

-- Step 1: Drop all time sub-partitions
-- This would need to be done for all 64 hash partitions
-- Example structure (needs careful implementation):

/*
DO $$
DECLARE
    partition_num INT;
    year_val INT;
    month_val INT;
    sub_partition_name TEXT;
BEGIN
    FOR partition_num IN 0..63 LOOP
        FOR year_val IN 2024..2025 LOOP
            FOR month_val IN 1..12 LOOP
                sub_partition_name := format('messages_p%s_%s_%02d', partition_num, year_val, month_val);
                EXECUTE format('DROP TABLE IF EXISTS %I CASCADE', sub_partition_name);
            END LOOP;
        END LOOP;
    END LOOP;
END $$;
*/

-- Step 2: Drop helper function
DROP FUNCTION IF EXISTS get_time_partition_name(TEXT, INT, INT);

-- IMPORTANT: This is a template
-- Actual rollback requires careful planning and data migration

