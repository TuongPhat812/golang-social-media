# Kế hoạch Partitioning và Migration

## 1. Cấu trúc hiện tại

### Table: `messages`
- `id` (UUID, PRIMARY KEY)
- `sender_id` (TEXT)
- `receiver_id` (TEXT)
- `content` (TEXT)
- `created_at` (TIMESTAMPTZ)

## 2. Cấu trúc mới với Partitioning

### Table: `messages` (Partitioned)
- `id` (UUID)
- `sender_id` (TEXT)
- `receiver_id` (TEXT)
- `content` (TEXT)
- `created_at` (TIMESTAMPTZ)
- `shard_id` (INT, GENERATED) - **Field mới để lưu kết quả hash**

### Partitions: 16 partitions (messages_p0 đến messages_p15)
- Partition key: `(sender_id, receiver_id)` - Hash partition
- `shard_id` được tính tự động: `hash(sender_id, receiver_id) % 16`

## 3. Lợi ích của `shard_id` field

### ✅ Dễ migrate khi tăng partitions:
- Khi tăng từ 16 → 32 partitions:
  - Dùng `shard_id` để tính partition mới: `shard_id % 32`
  - Không cần re-hash toàn bộ data
  - Chỉ cần redistribute data

### ✅ Query optimization:
- Có thể query trực tiếp partition: `WHERE shard_id = 5`
- Giúp debug và monitoring
- Có thể balance load dựa trên shard_id

### ✅ Consistency:
- Đảm bảo cùng một (sender_id, receiver_id) luôn có cùng shard_id
- Normalize: luôn dùng ID nhỏ hơn trước để consistent

## 4. Migration Plan

### Phase 1: Tạo migration (Migration 000002)
- Tạo function `calculate_shard_id()`
- Tạo partitioned table với `shard_id` field
- Tạo 16 partitions
- Migrate data từ table cũ
- Rename tables

### Phase 2: Update Application Code
- Update `MessageModel` để include `shard_id`
- Code không cần thay đổi logic (PostgreSQL tự route)
- `shard_id` là GENERATED column nên tự động tính

### Phase 3: Testing
- Test INSERT performance
- Test query performance
- Verify data consistency

## 5. Future Migration: Tăng số partitions (16 → 32)

### Khi cần scale lên 32 partitions:

```sql
-- Step 1: Tạo thêm 16 partitions mới (p16-p31)
CREATE TABLE messages_p16 PARTITION OF messages 
    FOR VALUES WITH (MODULUS 32, REMAINDER 16);
-- ... đến p31

-- Step 2: Redistribute data dựa trên shard_id
-- Data sẽ tự động redistribute khi INSERT mới
-- Hoặc có thể migrate data cũ:
UPDATE messages SET shard_id = shard_id % 32;
-- (PostgreSQL sẽ tự động move data đến partition mới)

-- Step 3: Update function để dùng 32 shards
CREATE OR REPLACE FUNCTION calculate_shard_id(
    sender_id TEXT, 
    receiver_id TEXT, 
    num_shards INT DEFAULT 32
) RETURNS INT AS $$ ... $$;
```

## 6. Code Changes Required

### Minimal changes:
- Update `MessageModel` struct (thêm `ShardID` field)
- GORM sẽ tự động handle GENERATED column
- Không cần thay đổi business logic

### Optional optimizations:
- Query by shard_id nếu cần
- Monitoring per-shard performance

