## Notification Service Architecture with Scylla Replication

```
apps/notification-service/
├── cmd/notification-service/
│   └── main.go
├── internal/
│   ├── application/
│   │   ├── command/
│   │   │   ├── contracts/
│   │   │   │   └── create_notification.command.contract.go
│   │   │   ├── create_notification.command.go     # use cases (welcome, event-driven)
│   │   │   └── dto/
│   │   │       └── create_notification.command.dto.go
│   │   ├── query/
│   │   │   ├── contracts/
│   │   │   │   └── list_notifications.query.contract.go
│   │   │   └── list_notifications.query.go
│   │   └── consumers/
│   │       └── user_created.consumer.go          # handles Kafka user.created event
│   ├── domain/
│   │   ├── notification/
│   │   │   └── entity.go                         # Notification entity + type enum
│   │   └── user/
│   │       └── entity.go                         # Replicated user profile (domain-level)
│   ├── infrastructure/
│   │   ├── eventbus/
│   │   │   ├── kafka_publisher.go                # publishes NotificationCreated
│   │   │   └── user_created_subscriber.go        # consumes user.created, triggers command
│   │   ├── persistence/
│   │   │   └── scylla/
│   │   │       ├── session.go                    # Scylla session bootstrap
│   │   │       ├── notification_repository.db.go # DB representation, queries
│   │   │       └── user_repository.db.go         # stores replicated users
│   │   └── http/
│   │       └── router.go                         # optional REST endpoints (list notifications)
│   └── interfaces/
│       └── rest/
│           ├── query/
│           │   ├── contracts/list_notifications.http.contract.go
│           │   └── list_notifications.http.handler.go
│           └── dto/
│               └── notification.http.response.go
└── ...
```

- Domain entity `notification.Notification` dùng chung cho mọi loại (field `Type`).
- Persistence Scylla có bảng `notifications_by_user` và `notification_users`, nhưng domain code chỉ thao tác qua repository.
- Kafka flow: `auth-service` → `UserCreated` → subscriber → `CreateNotificationCommand` (welcome) → insert Scylla → `NotificationCreated` (Kafka) → socket-service broadcast.

