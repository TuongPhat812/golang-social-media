# Chat Service User Replication

## Overview

Chat service giờ đã có khả năng replicate users từ auth-service thông qua Kafka events, giúp chat service hoạt động độc lập mà không cần query auth-service.

## Implementation

### 1. Database Migration

Đã tạo migration `000003_create_users_table` để tạo `users` table trong chat service database.

### 2. Components

- **UserModel**: Model cho user data
- **UserRepository**: Repository với methods: Upsert, FindByID, Exists
- **HandleUserCreatedCommand**: Command handler để xử lý UserCreated events
- **UserCreatedSubscriber**: Kafka consumer để subscribe `user.created` topic

### 3. Bootstrap

- Setup UserRepository
- Setup HandleUserCreatedCommand
- Setup UserCreatedSubscriber với consumer group ID

### 4. Main

- Start subscriber trong background (non-blocking)
- Cleanup subscriber khi shutdown

## Environment Variables

Cần thêm vào `.env.local` và `.env.local.docker`:

```bash
# Chat Service
CHAT_USER_GROUP_ID=chat-service-user
```

## Migration

Chạy migration để tạo users table:

```bash
make migration-up
```

## Flow

1. User register → Auth service publish `user.created` event
2. Chat service consume event → HandleUserCreatedCommand
3. UserRepository.Upsert() → Lưu user vào PostgreSQL
4. Chat service giờ có user data riêng → Hoạt động độc lập!

## Benefits

- ✅ Chat service không cần query auth-service
- ✅ Better performance (local data)
- ✅ Resilience (chat service vẫn hoạt động nếu auth-service down)
- ✅ Eventual consistency (data sync qua events)

## Next Steps

1. Add user validation trong CreateMessage command (optional)
2. Add user cache nếu cần (optional)
3. Handle user updates nếu auth-service có update user events

