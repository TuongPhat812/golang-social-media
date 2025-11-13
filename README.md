# Golang Social Media Prototype

This repository is organised as a mono-repo with four Go services under `apps/` (gateway, chat, notification, socket) and shared packages in `pkg/`. Each service keeps a Domain-Driven Design (DDD) layout so the domain, application, infrastructure, and interface layers remain isolated.

## Layout

```
.
├── apps/
│   ├── gateway/
│   ├── chat-service/
│   ├── notification-service/
│   └── socket-service/
└── pkg/               # shared config, events, generated protobuf code, etc.
```

- `apps/gateway`: Gin HTTP gateway orchestrating downstream calls.
- `apps/chat-service`: gRPC service creating chat messages, persisting them in Postgres, and emitting Kafka events.
- `apps/notification-service`: gRPC service consuming chat events, creating notifications, and emitting follow-up events.
- `apps/socket-service`: WebSocket service broadcasting events to connected clients.
- `pkg`: Reusable packages (`config`, `events`, protobuf stubs under `pkg/gen`, …).

## Project Layout

Every service exposes its entrypoint under `apps/<service>/cmd/<service>/main.go` and keeps implemention details under `apps/<service>/internal/`.

```
<service>
├── cmd/<service-name>/main.go     # Composition root / entrypoint
└── internal/
    ├── application/               # Application services and use-cases
    ├── domain/                    # (service-specific domain entities when added)
    ├── infrastructure/            # Framework and adapter implementations (gRPC, Kafka, etc.)
    └── interfaces/                # Transport layers (HTTP, gRPC, WebSocket, etc.)
```

Shared packages now live in `pkg/` so services can import `golang-social-media/pkg/<package>`.

### Protobuf Contracts

- Schemas live under `proto/`. Currently only `proto/chat/v1/chat_service.proto`.
- Generated Go code is committed under `pkg/gen/chat/v1`.
- Ensure you have `protoc` plus the Go plugins (`protoc-gen-go`, `protoc-gen-go-grpc`) on your `PATH`.
- Regenerate Go code with:

  ```bash
  make proto
  ```

## Docker Compose Setup

Two Compose files are provided:

- `docker-compose.infra.yml` — brings up Kafka, Postgres, Loki, Promtail, Kafka UI, and Grafana (no Zookeeper required thanks to KRaft mode).
- `docker-compose.app.yml` — runs the Go services inside containers (expects the infra stack to be running).

### Start Infrastructure

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.infra.yml up -d
```

This brings up Postgres, Kafka, and the observability stack on the shared `gsm-network`:

- Postgres: `localhost:5432` (host) / `gsm-postgres:5432` (containers)
- Kafka: `localhost:9092` (host) / `kafka:9092` (containers)
- Loki: `http://localhost:3100`
- Promtail scrapes `/var/log/app/*.log` and forwards to Loki.
- Grafana: `http://localhost:3000` (default `admin/admin`)
- Kafka UI: `http://localhost:8088`

### Start Application Services

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.app.yml up --build
```

Each service runs with code mounted from the host, using Go's `GOTOOLCHAIN=auto` to pull the correct toolchain. The containers connect to Kafka/Postgres and write logs to `/var/log/app` (mounted from `./logs`) so Promtail/Loki can ingest them.

Stop everything with:

```bash
docker compose -f docker-compose.app.yml down
docker compose -f docker-compose.infra.yml down
```

## Event Flow (Use Case: Create Chat Message)

1. Client calls `POST /chat/messages` on the `gateway`.
2. Gateway invokes `chat-service` over gRPC (`CreateMessage`).
3. `chat-service` creates a message aggregate and publishes a `ChatCreated` event to Kafka.
4. `notification-service` consumes the `ChatCreated` event, generates a notification, and publishes a `NotificationCreated` event.
5. `socket-service` listens for both `ChatCreated` and `NotificationCreated` events and pushes real-time updates to connected clients.

## Manual Run (Without Docker)

If you prefer running binaries directly on the host, the executables automatically load environment variables from the repository’s `.env` file. Simply start each service in separate terminals:

```bash
# In separate shells
cd apps/notification-service && go run ./cmd/notification-service
cd apps/chat-service && go run ./cmd/chat-service
cd apps/socket-service && go run ./cmd/socket-service
cd apps/gateway && go run ./cmd/gateway
```

Override any setting by exporting it before launch, for example `export KAFKA_BROKERS=localhost:9094`.

## Database & Migrations

- Chat-service connects to Postgres using the `CHAT_DATABASE_DSN` environment variable (defaults to `postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable`).
- Schema changes rely on explicit SQL migrations (GORM auto-migrate is disabled).
- Versioned SQL migrations live in `apps/chat-service/migrations`. You can either use the Makefile targets or run the built-in Go CLI. Examples (from the repo root):

  ```bash
  cd apps/chat-service
  go run ./cmd/migrate create add_messages_index   # generates stub files
  go run ./cmd/migrate up                          # applies pending migrations
  go run ./cmd/migrate down                        # rolls back one step
  ```

  If you prefer Makefile aliases from the repo root you can still use:

  ```bash
  # create a new migration (NAME is required)
  make migration-create NAME=add_new_table

  # apply migrations (uses CHAT_DB_DSN or defaults to local postgres)
  make migration-up

  # rollback one step
  make migration-down
  ```

  Then edit the generated `.up.sql` / `.down.sql` files.

## Development Notes

- The repository uses a Go workspace (`go.work`) that includes every service under `apps/` plus the shared `pkg/` module.
- Kafka producers/consumers are implemented with [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go). When running outside Docker, use the host listener `localhost:9094`; services inside Docker should continue to use `kafka:9092`.
- Shared code is pulled from `pkg/` via local replace directives, so nothing needs to be published externally.
- Logging uses Zerolog with JSON output. Set `LOG_OUTPUT_DIR` (default `./logs`) to control file location. Docker compose sets it to `/var/log/app`, which Promtail watches and ships to Loki/Grafana.
- Gateway-specific environment flags:
  - `GIN_MODE=release` (or `debug`/`test`) to control Gin’s mode.
  - `GIN_DISABLE_ACCESS_LOG=true` hides Gin’s access log noise (e.g. per-request HTTP logs). Leave unset or `false` to keep the default output.

Future work will flesh out persistence, authentication, and real-time delivery handlers while keeping the DDD boundaries intact.
