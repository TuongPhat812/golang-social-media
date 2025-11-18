# Environment Variables

All services automatically load environment variables from the repository's `.env` file (via `pkg/config`).

## Common Variables

- `KAFKA_BROKERS`: Kafka broker addresses (comma-separated). Default: `localhost:9092` (host) / `kafka:9092` (Docker)
- `LOG_OUTPUT_DIR`: Directory for log files. Default: `./logs`. Docker sets it to `/var/log/app`

## Service-Specific Variables

### Gateway

- `GIN_MODE`: Gin's mode (`release`, `debug`, or `test`). Default: `debug`
- `GIN_DISABLE_ACCESS_LOG`: Set to `true` to hide Gin's access log noise (e.g. per-request HTTP logs). Leave unset or `false` to keep the default output
- `AUTH_SERVICE_URL`: Base URL for the auth service REST API. Default: `http://localhost:9101`
- `CHAT_SERVICE_ADDR`: gRPC address for chat service. Default: `localhost:9000`

### Chat Service

- `CHAT_DATABASE_DSN`: PostgreSQL connection string. Default: `postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable`
- `CHAT_SERVICE_PORT`: gRPC server port. Default: `9000`

### Notification Service

- `SCYLLA_HOSTS`: ScyllaDB hosts (comma-separated). Default: `localhost:9042` (single-node) or `localhost:9042,localhost:9043,localhost:9044` (3-node cluster)
- `SCYLLA_KEYSPACE`: ScyllaDB keyspace name. Default: `notification_service`
- `NOTIFICATION_USER_GROUP_ID`: Kafka consumer group ID for `user.created` events. Default: `notification-service-user`
- `NOTIFICATION_CHAT_GROUP_ID`: Kafka consumer group ID for `chat.created` events. Default: `notification-service-chat`

### Socket Service

- `SOCKET_CHAT_GROUP_ID`: Kafka consumer group ID for `chat.created` events. Default: `socket-service-chat`
- `SOCKET_NOTIFICATION_GROUP_ID`: Kafka consumer group ID for `notification.created` events. Default: `socket-service-notification`
- `SOCKET_SERVICE_PORT`: WebSocket server port. Default: `9200`

## Override at Runtime

Override any setting by exporting it before launch:

```bash
export KAFKA_BROKERS=localhost:9094
export CHAT_DATABASE_DSN=postgres://user:pass@localhost:5432/db
```

