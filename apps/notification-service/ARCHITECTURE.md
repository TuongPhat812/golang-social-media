## Notification Service Architecture with Scylla

```
apps/notification-service/
├── cmd/notification-service/
│   └── main.go                    # composition root
├── internal/
│   ├── application/
│   │   ├── command/
│   │   │   ├── contracts/create_notification.command.contract.go
│   │   │   ├── create_notification.command.go
│   │   │   └── dto/create_notification.command.dto.go
│   │   ├── query/
│   │   │   ├── contracts/list_notifications.query.contract.go
│   │   │   └── list_notifications.query.go
│   │   └── consumers/
│   │       └── user_created.consumer.go        # consumes user.created events
│   ├── domain/
│   │   ├── notification/entity.go
│   │   └── user/entity.go
│   ├── infrastructure/
│   │   ├── eventbus/
│   │   │   ├── kafka_publisher.go             # emits NotificationCreated
│   │   │   └── user_created_subscriber.go
│   │   ├── persistence/
│   │   │   └── scylla/
│   │   │       ├── session.go
│   │   │       ├── notification_repository.db.go
│   │   │       └── user_repository.db.go
│   │   └── http/router.go (optional query API)
│   └── interfaces/
│       └── rest/query/
│           ├── contracts/list_notifications.http.contract.go
│           └── list_notifications.http.handler.go
└── ...
```

### Flow
1. Auth-service publishes `UserCreated` → consumer `user_created.consumer.go`.
2. Consumer replicates user via `user_repository.db.go` and triggers `CreateNotificationCommand`.
3. Command writes `Notification` domain object to Scylla (`notifications_by_user` table) using repository.
4. Publisher emits `NotificationCreated` (Kafka) → socket-service broadcasts.
5. Optional REST query returns notifications using unified domain entity; repository handles DB specifics.

