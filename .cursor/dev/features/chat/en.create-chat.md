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
   - Loads env, sets up Kafka publisher, registers gRPC handler.
5. **gRPC handler** – `apps/chat-service/internal/interfaces/grpc/chat/handler.go`
   - Receives protobuf request, forwards to application service.
6. **Use case / Application service** – `apps/chat-service/internal/application/messages/service.go`
   - Builds domain `Message`, delegates persistence/side effects via the `EventPublisher` interface.

==============================

## 2. Side-Effect Logic (Kafka & Consumers)
7. **Kafka publisher** – `apps/chat-service/internal/infrastructure/eventbus/kafka_publisher.go`
   - Implements `EventPublisher`; emits `events.ChatCreated` to topic `chat.created`.
8. **Notification consumer** – `apps/notification-service/internal/infrastructure/eventbus/subscriber.go`
   - Consumes `chat.created`, dispatches to application service.
9. **Notification application** – `apps/notification-service/internal/application/notifications/service.go`
   - Builds notification domain object, emits `events.NotificationCreated`.
10. **Notification publisher** – `apps/notification-service/internal/infrastructure/eventbus/kafka_publisher.go`
    - Publishes `notification.created` for downstream listeners.
11. **Socket listeners** – `apps/socket-service/internal/infrastructure/eventbus/listener.go`
    - Two consumers: `chat.created` and `notification.created`.
12. **Socket application service** – `apps/socket-service/internal/application/events/service.go`
    - Logs and forwards events to the WebSocket hub.
13. **WebSocket hub** – `apps/socket-service/internal/interfaces/socket/hub.go`
    - Currently logs broadcasts; extend to push to connected clients.

## 3. Shared Contracts & Configuration
14. **Protobuf contract** – `proto/chat/v1/chat_service.proto` + generated Go code in `pkg/gen/chat/v1`.
15. **Event payloads** – `pkg/events/chat.go`, `pkg/events/notification.go`.
16. **Environment loader** – `pkg/config/env.go` (`KAFKA_BROKERS`, `<SERVICE>_PORT`, etc.).
17. **Docker Compose** – `docker-compose.infra.yml`, `docker-compose.app.yml` (brokers: `kafka:9092` vs `localhost:9094`).
