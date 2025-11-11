# Golang Social Media Prototype

This repository contains five Go modules representing a gateway HTTP API, two gRPC microservices, a dedicated WebSocket server, and shared domain entities. Each service follows a Domain-Driven Design (DDD) layout to keep the domain, application, infrastructure, and interface layers separated for future growth.

## Modules

- `gateway`: Gin-based HTTP gateway serving REST endpoints and orchestrating downstream calls to backend services.
- `chat-service`: gRPC server responsible for chat-related features and emitting chat domain events to Kafka.
- `notification-service`: gRPC server that reacts to chat events, creates notifications, and emits notification events to Kafka.
- `socket-service`: WebSocket gateway that listens for domain events from Kafka and broadcasts payloads to connected clients.
- `common`: Shared domain entities, contracts, and event definitions consumed across services.

## Project Layout

Every service module exposes the executable inside `cmd/<service-name>` and keeps internal code within the `internal` directory.

```
<service>
├── cmd/<service-name>/main.go     # Composition root / entrypoint
└── internal/
    ├── application/               # Application services and use-cases
    ├── domain/                    # (service-specific domain entities when added)
    ├── infrastructure/            # Framework and adapter implementations (gRPC, Kafka, etc.)
    └── interfaces/                # Transport layers (HTTP, gRPC, WebSocket, etc.)
```

The `common` module stores shared domain types, contracts, and event payloads so services can rely on a single source of truth.

## Docker Compose Setup

Two Compose files are provided:

- `docker-compose.infra.yml` — brings up a single-node Kafka broker using the official `apache/kafka` KRaft image (no Zookeeper required).
- `docker-compose.app.yml` — runs the Go services inside containers (expects the infra stack to be running).

### Start Infrastructure

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.infra.yml up -d
```

This creates the shared `gsm-network` and exposes Kafka on `kafka:9092` for other containers and `localhost:9092` for host clients.

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
cd notification-service && go run ./cmd/notification-service
cd chat-service && go run ./cmd/chat-service
cd socket-service && go run ./cmd/socket-service
cd gateway && go run ./cmd/gateway
```

Override any setting by exporting it before launch, for example `export KAFKA_BROKERS=localhost:9092`.

## Development Notes

- The repository uses a Go workspace (`go.work`) so modules can reference each other locally.
- Kafka producers/consumers are implemented with [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go). When running outside Docker, use the host listener `localhost:9092`; services inside Docker should continue to use `kafka:9092`.
- Dependencies on the `common` module are resolved via local replace directives and do not require publishing to a remote repository.

Future work will flesh out persistence, authentication, and real-time delivery handlers while keeping the DDD boundaries intact.
