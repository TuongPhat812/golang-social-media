# Event Flow

## Use Case: Create Chat Message

1. Client calls `POST /chat/messages` on the `gateway`.
2. Gateway invokes `chat-service` over gRPC (`CreateMessage`).
3. `chat-service` creates a message aggregate and publishes a `ChatCreated` event to Kafka.
4. `notification-service` consumes the `ChatCreated` event, generates a notification, and publishes a `NotificationCreated` event.
5. `socket-service` listens for both `ChatCreated` and `NotificationCreated` events and pushes real-time updates to connected clients.

## Use Case: User Registration

1. Client calls `POST /auth/register` on the `gateway`.
2. Gateway invokes `auth-service` over HTTP REST API (`RegisterUser`).
3. `auth-service` creates a user and publishes a `UserCreated` event to Kafka.
4. `notification-service` consumes the `UserCreated` event:
   - Replicates user data into ScyllaDB
   - Creates a welcome notification
   - Publishes a `NotificationCreated` event
5. `socket-service` consumes the `NotificationCreated` event and broadcasts it to connected clients via WebSocket.

## Event Topics

Events are published to the following Kafka topics:

- `user.created` - Published when a new user is registered
- `chat.created` - Published when a new chat message is created
- `notification.created` - Published when a new notification is created

## Event Payloads

Event payloads are defined in `pkg/events/`:

- `pkg/events/user.go` - `UserCreated` event
- `pkg/events/chat.go` - `ChatCreated` event
- `pkg/events/notification.go` - `NotificationCreated` event
- `pkg/events/topics.go` - Topic name constants

## Consumer Groups

Each service uses specific consumer group IDs to ensure proper event distribution:

- `notification-service-user` - Consumes `user.created` events
- `notification-service-chat` - Consumes `chat.created` events
- `socket-service-chat` - Consumes `chat.created` events
- `socket-service-notification` - Consumes `notification.created` events

These can be configured via environment variables (see [environment.md](./environment.md)).

