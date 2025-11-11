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
└── pkg/               # shared config, contracts, events, codecs, etc.
```

- `apps/gateway`: Gin HTTP gateway orchestrating downstream calls.
- `apps/chat-service`: gRPC service creating chat messages and emitting Kafka events.
- `apps/notification-service`: gRPC service consuming chat events, creating notifications, and emitting follow-up events.
- `apps/socket-service`: WebSocket service broadcasting events to connected clients.
- `pkg`: Reusable packages (`config`, `contracts`, `events`, `grpcjson`, …).

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

Shared packages now live in `pkg/` so services can import `github.com/<org>/golang-social-media/pkg/<package>`.

## Docker Compose Setup

Two Compose files are provided:

- `docker-compose.infra.yml` — brings up a single-node Kafka broker using the official `apache/kafka` KRaft image (no Zookeeper required).
- `docker-compose.app.yml` — runs the Go services inside containers (expects the infra stack to be running).

### Start Infrastructure

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.infra.yml up -d
```

This creates the shared `gsm-network` and exposes Kafka on `kafka:9092` for other containers and `localhost:9094` for host clients.

### Start Application Services

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.app.yml up --build
```

Each service runs with code mounted from the host, using Go's `GOTOOLCHAIN=auto` to pull the correct toolchain. The containers connect to Kafka through the shared Docker network.

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

## Development Notes

- The repository uses a Go workspace (`go.work`) that includes every service under `apps/` plus the shared `pkg/` module.
- Kafka producers/consumers are implemented with [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go). When running outside Docker, use the host listener `localhost:9094`; services inside Docker should continue to use `kafka:9092`.
- Shared code is pulled from `pkg/` via local replace directives, so nothing needs to be published externally.

Future work will flesh out persistence, authentication, and real-time delivery handlers while keeping the DDD boundaries intact.
