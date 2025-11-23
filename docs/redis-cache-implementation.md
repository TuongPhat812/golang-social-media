# Redis Cache Implementation

## Overview

Đã implement Redis cache cho các services để cải thiện performance và giảm tải database.

## Infrastructure

### Redis Service
- **Image**: `redis:7-alpine`
- **Port**: `6379`
- **Persistence**: AOF (Append Only File) enabled
- **Healthcheck**: Redis ping
- **Volume**: `redis-data` for persistence

## Shared Cache Package

### `pkg/cache/cache.go`
- **Cache Interface**: Standard cache operations (Get, Set, Delete, DeletePattern, Exists, Close)
- **RedisCache Implementation**: Redis-backed cache với connection pooling
- **Features**:
  - Connection pooling (PoolSize: 10, MinIdleConns: 5)
  - Retry mechanism (MaxRetries: 3)
  - Timeouts (Dial: 5s, Read: 3s, Write: 3s)
  - Error handling với logging

## Chat Service Cache

### User Cache (`apps/chat-service/internal/infrastructure/cache/user.cache.go`)
- **Purpose**: Cache user data replicated from auth-service
- **TTL**: 30 minutes
- **Key Pattern**: `chat:user:{user_id}`
- **Methods**:
  - `GetUser(ctx, id)`: Get user from cache
  - `SetUser(ctx, user)`: Store user in cache
  - `DeleteUser(ctx, id)`: Remove user from cache

### Integration
- **UserRepository**: Cache-aside pattern
  - `FindByID()`: Check cache first, fallback to DB, update cache
  - `Upsert()`: Save to DB, update cache
  - `Exists()`: Check cache first, fallback to DB

## Auth Service Cache

### User Cache (`apps/auth-service/internal/infrastructure/cache/user.cache.go`)
- **Purpose**: Cache user profiles for faster lookups
- **TTL**: 15 minutes
- **Key Patterns**:
  - `auth:user:id:{user_id}`: Cache by user ID
  - `auth:user:email:{email}`: Cache by email
- **Methods**:
  - `GetUser(ctx, id)`: Get user by ID
  - `SetUser(ctx, user)`: Store user by ID
  - `GetUserByEmail(ctx, email)`: Get user by email
  - `SetUserByEmail(ctx, user)`: Store user by email
  - `DeleteUser(ctx, id)`: Remove user from cache

### Integration
- **UserRepository**: In-memory repository (có thể wrap với cache sau)
- Cache có thể dùng cho future database migration

## Environment Variables

### All Services:
```bash
REDIS_ADDR=localhost:6379        # Redis address
REDIS_PASSWORD=                   # Redis password (empty for local)
REDIS_DB=0                        # Redis database number
```

### Docker:
```bash
REDIS_ADDR=redis:6379             # Use container name in Docker
```

## Cache Patterns

### Cache-Aside (Lazy Loading)
1. Check cache
2. If cache miss → Query database
3. Store result in cache
4. Return result

### Write-Through
1. Write to database
2. Update cache
3. Return success

## Benefits

- ✅ **Performance**: Faster reads (Redis in-memory)
- ✅ **Reduced DB Load**: Less queries to PostgreSQL
- ✅ **Scalability**: Can handle more concurrent requests
- ✅ **Resilience**: Service continues if cache fails (graceful degradation)

## Cache Invalidation

### Automatic:
- TTL-based expiration (15-30 minutes)
- Manual deletion on updates

### Manual:
- `DeleteUser()`: Remove specific user
- `DeletePattern()`: Remove multiple keys (e.g., `chat:user:*`)

## Next Steps

1. Add cache metrics (hit/miss rates)
2. Implement cache warming strategies
3. Add cache invalidation on user updates
4. Consider cache clustering for high availability

