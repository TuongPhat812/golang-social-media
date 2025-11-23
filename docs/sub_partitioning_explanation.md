# Sub-Partitioning (Nested Partitioning) theo thời gian

## Câu hỏi: 1 partition có thể chia nhỏ thêm theo thời gian không?

**Trả lời: CÓ!** PostgreSQL hỗ trợ **Sub-Partitioning** (hay Nested Partitioning) từ version 10+.

## Cách hoạt động:

### Cấu trúc 2-level partitioning:

```
messages (parent table)
├── Hash Partition by (sender_id, receiver_id) - 64 partitions
│   ├── messages_p0
│   │   ├── Range Partition by created_at (monthly)
│   │   │   ├── messages_p0_2024_01
│   │   │   ├── messages_p0_2024_02
│   │   │   └── messages_p0_2024_03
│   ├── messages_p1
│   │   ├── messages_p1_2024_01
│   │   ├── messages_p1_2024_02
│   │   └── messages_p1_2024_03
│   └── ... (64 hash partitions, mỗi partition có nhiều time sub-partitions)
```

## Lợi ích của Sub-Partitioning:

### ✅ Level 1: Hash Partition (64 partitions)
- Giảm lock contention: 64x ít lock hơn
- Parallel writes: Nhiều partitions write đồng thời
- Load balancing: Distribute writes đều

### ✅ Level 2: Range Partition theo thời gian (monthly/quarterly)
- **Dễ dàng drop old data**: `DROP TABLE messages_p0_2023_01` (xóa data cũ)
- **Smaller indexes**: Mỗi time partition có index nhỏ hơn
- **Better query performance**: Partition pruning theo cả hash và time
- **Maintenance**: Dễ dàng archive/delete data cũ

## Ví dụ Implementation:

```sql
-- Level 1: Hash partition (64 partitions)
CREATE TABLE messages (
    ...
) PARTITION BY HASH (sender_id, receiver_id);

-- Level 2: Mỗi hash partition lại partition theo thời gian
CREATE TABLE messages_p0 PARTITION OF messages
    FOR VALUES WITH (MODULUS 64, REMAINDER 0)
    PARTITION BY RANGE (created_at);

-- Sub-partitions cho messages_p0
CREATE TABLE messages_p0_2024_01 PARTITION OF messages_p0
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE messages_p0_2024_02 PARTITION OF messages_p0
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
-- ... tiếp tục cho các tháng khác
```

## Khi nào nên dùng Sub-Partitioning?

### ✅ Nên dùng khi:
- Data có **time-based access pattern** (query theo thời gian)
- Cần **retention policy** (xóa data cũ định kỳ)
- Table rất lớn (hàng trăm GB đến TB)
- Cần **archive old data** dễ dàng

### ❌ Không cần khi:
- Data nhỏ (< 100GB)
- Không có nhu cầu xóa data cũ
- Query không theo pattern thời gian

## Trade-offs:

### Pros:
- ✅ Maximum performance: 64 hash partitions × N time partitions
- ✅ Easy data lifecycle management
- ✅ Better query performance với partition pruning

### Cons:
- ❌ Phức tạp hơn: Phải maintain nhiều partitions
- ❌ Overhead: Nhiều partitions = nhiều metadata
- ❌ Migration phức tạp hơn

## Recommendation cho project này:

### Hiện tại (64 hash partitions):
- ✅ Đủ tốt cho 6k-20k req/s
- ✅ Đơn giản, dễ maintain
- ✅ Performance tốt

### Nếu cần Sub-Partitioning sau này:
- Khi data > 100GB
- Khi cần retention policy (xóa data > 1 năm)
- Khi query pattern theo thời gian nhiều

## Migration path:

Nếu muốn thêm sub-partitioning sau này:
1. Tạo migration mới để convert hash partitions thành sub-partitioned
2. Migrate data từ hash partition → hash + time partition
3. Update application code nếu cần (thường không cần)

