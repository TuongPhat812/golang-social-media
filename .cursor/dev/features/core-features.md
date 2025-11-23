# Core Features Summary

Tổng hợp các tính năng đã được implement trong project.

## 🏗️ Architecture Overview

- **Microservices Architecture**: 5 services (gateway, auth-service, chat-service, notification-service, socket-service)
- **Domain-Driven Design (DDD)**: Mỗi service follow DDD pattern với domain, application, infrastructure, interfaces layers
- **Event-Driven Architecture**: Kafka-based event streaming giữa các services
- **CQRS Pattern**: Tách biệt Command và Query trong gateway và các services
- **Observability**: Grafana, Loki, Prometheus cho logging và metrics

---

## 🔐 Auth Service (`apps/auth-service`)

### Features Implemented:
1. **User Registration** (`POST /auth/register`)
   - Tạo user mới với validation (email, password, name)
   - Publish `UserCreated` event lên Kafka topic `user.created`
   - In-memory storage (có thể migrate sang database sau)

2. **User Login** (`POST /auth/login`)
   - Xác thực email/password
   - Generate token (simple token store, có thể upgrade JWT sau)
   - Return user ID và token

3. **Get User Profile** (`GET /auth/profile/:id`)
   - Lấy thông tin user theo ID
   - Return user profile (ID, email, name)

### Domain Events:
- `UserCreated` - Published khi user đăng ký thành công

### Technology:
- **Transport**: HTTP REST (Gin)
- **Storage**: In-memory (có thể migrate sang PostgreSQL)
- **Event Bus**: Kafka publisher

---

## 💬 Chat Service (`apps/chat-service`)

### Features Implemented:
1. **Create Message** (`CreateMessage` gRPC method)
   - Tạo chat message giữa 2 users
   - Persist vào PostgreSQL với **64 hash partitions** (optimized for high concurrency)
   - Publish `ChatCreated` event lên Kafka topic `chat.created`
   - Auto-calculate `shard_id` từ (sender_id, receiver_id) để support future partition scaling

### Database:
- **PostgreSQL** với partitioning strategy:
  - **64 hash partitions** (messages_p0 đến messages_p63)
  - Partition key: `HASH (sender_id, receiver_id)`
  - `shard_id` field: GENERATED column để dễ migrate khi tăng partitions
  - Indexes trên mỗi partition: sender_id, receiver_id, created_at, shard_id

### Performance Optimizations:
- Connection pooling (MaxOpenConns: 100, MaxIdleConns: 25)
- gRPC server optimizations (MaxConcurrentStreams: 10000, increased window sizes)
- PostgreSQL WAL tuning (wal_buffers, checkpoint_timeout, max_wal_size, etc.)

### Domain Events:
- `MessageCreated` - Published khi message được tạo thành công

### Technology:
- **Transport**: gRPC
- **Storage**: PostgreSQL (GORM)
- **Event Bus**: Kafka publisher

---

## 🔔 Notification Service (`apps/notification-service`)

### Features Implemented:
1. **Event Consumers**:
   - Consume `user.created` events → Replicate user data vào ScyllaDB
   - Consume `chat.created` events → Tạo notification cho receiver
   - Publish `NotificationCreated` event lên Kafka topic `notification.created`

2. **Notification Management**:
   - Store notifications trong ScyllaDB (notification_service keyspace)
   - Query notifications by user ID
   - Mark notifications as read

### Database:
- **ScyllaDB** (3-node cluster):
  - Keyspace: `notification_service`
  - Tables: `notifications_by_user` (partitioned by user_id)
  - Consistency level: QUORUM (configurable via `SCYLLA_CONSISTENCY_LEVEL`)

### Kafka Consumers:
- `notification-service-user` group: Consume `user.created`
- `notification-service-chat` group: Consume `chat.created`
- Optimized với timeouts, retries, batching

### Domain Events:
- `NotificationCreated` - Published khi notification được tạo

### Technology:
- **Transport**: gRPC (for future queries)
- **Storage**: ScyllaDB (gocql driver)
- **Event Bus**: Kafka consumer + publisher

---

## 🔌 Socket Service (`apps/socket-service`)

### Features Implemented:
1. **WebSocket Server**:
   - WebSocket endpoint để clients connect
   - Broadcast real-time events đến connected clients

2. **Event Consumers**:
   - Consume `chat.created` events → Broadcast đến WebSocket clients
   - Consume `notification.created` events → Broadcast đến WebSocket clients

### Kafka Consumers:
- `socket-service-chat` group: Consume `chat.created`
- `socket-service-notification` group: Consume `notification.created`
- Optimized với timeouts, retries, batching

### Technology:
- **Transport**: WebSocket (Gorilla WebSocket)
- **Event Bus**: Kafka consumer

---

## 🌐 Gateway Service (`apps/gateway`)

### Features Implemented:
1. **HTTP Endpoints**:
   - `POST /chat/messages` - Create chat message (delegates to chat-service via gRPC)
   - `POST /auth/register` - Register user (delegates to auth-service via HTTP)
   - `POST /auth/login` - Login user (delegates to auth-service via HTTP)
   - `GET /auth/profile/:id` - Get user profile (delegates to auth-service via HTTP)

### Architecture:
- **API Gateway Pattern**: Single entry point cho tất cả clients
- **Service Orchestration**: Orchestrate calls đến downstream services
- **CQRS**: Tách biệt Command handlers và Query handlers
- **Error Handling**: Centralized error handling với error transformer

### Technology:
- **Transport**: HTTP REST (Gin)
- **Downstream Services**: gRPC (chat-service), HTTP (auth-service)
- **Error Handling**: Custom error system với error codes

---

## 📊 Infrastructure & Observability

### Infrastructure Services:
1. **PostgreSQL** (chat-service database)
   - Optimized configuration cho write performance
   - 64 hash partitions cho messages table
   - Connection pooling

2. **ScyllaDB** (notification-service database)
   - 3-node cluster
   - QUORUM consistency level
   - Optimized connection settings

3. **Kafka** (event streaming)
   - Topics: `user.created`, `chat.created`, `notification.created`
   - Consumer groups cho mỗi service

4. **Observability Stack**:
   - **Loki**: Log aggregation
   - **Promtail**: Log shipper
   - **Prometheus**: Metrics collection
   - **Grafana**: Visualization
   - **Kafka UI**: Kafka management UI
   - **Cassandra Web UI**: ScyllaDB management UI

---

## 🔄 Event Flow

### User Registration Flow:
1. Client → Gateway: `POST /auth/register`
2. Gateway → Auth Service: HTTP call
3. Auth Service: Create user, publish `UserCreated` event
4. Notification Service: Consume `UserCreated`, replicate to ScyllaDB, create welcome notification
5. Socket Service: Consume `NotificationCreated`, broadcast via WebSocket

### Chat Message Flow:
1. Client → Gateway: `POST /chat/messages`
2. Gateway → Chat Service: gRPC `CreateMessage`
3. Chat Service: Persist to PostgreSQL (64 partitions), publish `ChatCreated` event
4. Notification Service: Consume `ChatCreated`, create notification, publish `NotificationCreated`
5. Socket Service: Consume `ChatCreated` và `NotificationCreated`, broadcast via WebSocket

---

## 🚀 Performance Features

### Database Optimizations:
- **PostgreSQL**:
  - 64 hash partitions (giảm lock contention)
  - Connection pooling
  - WAL tuning (wal_buffers, checkpoint_timeout, max_wal_size)
  - Indexes trên mỗi partition

- **ScyllaDB**:
  - 3-node cluster với QUORUM consistency
  - Optimized connection timeouts và retry policies

### Application Optimizations:
- **gRPC**: High concurrency settings (MaxConcurrentStreams: 10000)
- **Kafka**: Optimized consumer settings (timeouts, batching, retries)
- **Load Testing**: Script `scripts/load_test_chat.go` để test performance

### Expected Performance:
- **Native**: ~6,000-6,700 req/s
- **Docker**: ~800 req/s (Mac Docker Desktop overhead)
- **With 64 partitions**: Expected 20-40k req/s (theoretical)

---

## 📝 Database Migrations

### PostgreSQL (chat-service):
- Migration system: `golang-migrate/migrate/v4`
- Migrations: `apps/chat-service/migrations/`
- Commands:
  - `make migration-create NAME=<name>` - Create new migration
  - `make migration-up` - Apply migrations
  - `make migration-down` - Rollback one migration

### ScyllaDB (notification-service):
- CQL scripts: `apps/notification-service/infra/scylla/`
- Apply via: `docker exec gsm-scylla-1 cqlsh -f /var/lib/scylla-init/notification_service.cql`

---

## 🛠️ Development Tools

### Makefile Commands:
- `make proto` - Generate protobuf Go code
- `make migration-create NAME=<name>` - Create migration
- `make migration-up` - Apply migrations
- `make migration-down` - Rollback migration
- `make load-test-chat` - Run load test

### Scripts:
- `reset_infra.sh` - Reset infrastructure (down, cleanup, up, migrate)
- `scripts/load_test_chat.go` - Load testing script

### Environment Files:
- `.env.local` - Local development (host execution)
- `.env.local.docker` - Docker Compose execution
- `SKIP_ENV_FILE=true` - Skip loading .env in Docker containers

---

## 📚 Documentation Structure

- `.cursor/dev/` - Detailed documentation
  - `setup.md` - Development environment setup
  - `docker.md` - Docker Compose guide
  - `migrations.md` - Database migrations
  - `environment.md` - Environment variables
  - `running.md` - Running services
  - `event-flow.md` - Event-driven architecture
  - `source-guide.md` - Codebase navigation
  - `features/` - Feature-specific documentation

---

## 🔮 Future Enhancements (Not Yet Implemented)

### Potential Features:
- JWT authentication (hiện tại dùng simple token)
- Friend requests/acceptance
- Chat thread listing
- Message history queries
- User search
- Real-time presence
- File uploads
- Message reactions
- Read receipts

### Database Improvements:
- Sub-partitioning theo thời gian (nếu cần retention policy)
- Read replicas cho PostgreSQL
- ScyllaDB multi-datacenter setup

### Infrastructure:
- Service mesh (Istio/Linkerd)
- API rate limiting
- Circuit breakers
- Distributed tracing (Jaeger/Zipkin)

---

## 📊 Current Status

✅ **Implemented & Working**:
- User registration/login/profile
- Chat message creation
- Event-driven notifications
- WebSocket real-time updates
- 64-partition PostgreSQL setup
- ScyllaDB 3-node cluster
- Full observability stack
- Load testing tools

🚧 **In Progress / Optimized**:
- Performance tuning (64 partitions)
- Kafka consumer optimizations
- Database connection pooling

📋 **Planned**:
- Additional features (see Future Enhancements)
- Production-ready improvements

