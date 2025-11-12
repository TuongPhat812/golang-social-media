# Feature Walkthrough: Create Chat Message

Use this checklist to trace `POST /chat/messages` end-to-end.

## 1. Main Feature Flow
1. **Route definition** – `apps/gateway/internal/interfaces/rest/handlers.go`
   - Handler validates `createMessageRequest`, then calls the application service.
2. **Application service** – `apps/gateway/internal/application/messages/service.go`
   - Implements the `Service` interface. Invokes `ChatClient.CreateMessage` and maps the protobuf response to the gateway domain entity.
3. **gRPC client** – `apps/gateway/internal/infrastructure/grpc/chat/client.go`
   - Creates `pkg/gen/chat/v1.ChatServiceClient`, loads `CHAT_SERVICE_ADDR` from `pkg/config`.
4. **Server bootstrap** – `apps/chat-service/cmd/chat-service/main.go`
   - Loads env, opens a GORM connection to Postgres (`CHAT_DATABASE_DSN`), sets up Kafka publisher, and registers gRPC handler. (Migrations are manual via `go run ./cmd/migrate`.)
5. **gRPC handler** – `apps/chat-service/internal/interfaces/grpc/chat/handler.go`
   - Receives protobuf request, forwards to application service.
6. **Use case / Application service** – `apps/chat-service/internal/application/messages/service.go`
   - Builds domain `Message`, persists it through the repository, then raises the domain event.

==============================

## 2. Persistence & Side-Effect Logic
7. **Repository implementation** – `apps/chat-service/internal/infrastructure/persistence/message_repository.go`
   - Wraps GORM and satisfies the `messages.Repository` interface.

8. **Kafka publisher** – `apps/chat-service/internal/infrastructure/eventbus/kafka_publisher.go`
   - Implements `EventPublisher`; emits `events.ChatCreated` to topic `chat.created`.
9. **Notification consumer** – `apps/notification-service/internal/infrastructure/eventbus/subscriber.go`
   - Consumes `chat.created`, dispatches to application service.
10. **Notification application** – `apps/notification-service/internal/application/notifications/service.go`
   - Builds notification domain object, emits `events.NotificationCreated`.
11. **Notification publisher** – `apps/notification-service/internal/infrastructure/eventbus/kafka_publisher.go`
    - Publishes `notification.created` for downstream listeners.
12. **Socket listeners** – `apps/socket-service/internal/infrastructure/eventbus/listener.go`
    - Two consumers: `chat.created` and `notification.created`.
13. **Socket application service** – `apps/socket-service/internal/application/events/service.go`
    - Logs and forwards events to the WebSocket hub.
14. **WebSocket hub** – `apps/socket-service/internal/interfaces/socket/hub.go`
    - Currently logs broadcasts; extend to push to connected clients.

## 3. Shared Contracts & Configuration
15. **Protobuf contract** – `proto/chat/v1/chat_service.proto` + generated Go code in `pkg/gen/chat/v1`.
16. **Event payloads** – `pkg/events/chat.go`, `pkg/events/notification.go`.
17. **Environment loader** – `pkg/config/env.go` (`KAFKA_BROKERS`, `<SERVICE>_PORT`, `CHAT_DATABASE_DSN`, etc.).
18. **Docker Compose** – `docker-compose.infra.yml`, `docker-compose.app.yml` (brokers: `kafka:9092` vs `localhost:9094`, Postgres: `gsm-postgres:5432` vs `localhost:5432`).
