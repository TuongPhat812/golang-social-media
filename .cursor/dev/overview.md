# Project At A Glance

**Name:** Golang Social Media Prototype

**Purpose:** Provide a sample social-media backend composed of multiple Go services (gateway, chat, notification, socket) communicating via gRPC and Kafka. The repo demonstrates a DDD-inspired layout inside a mono-repo.

## Core Concepts

- **apps/** — contains all runnable services. Each service follows a `cmd/` entrypoint plus `internal/` layered packages (domain, application, infrastructure, interfaces).
- **pkg/** — reusable building blocks shared across services (`config`, `events`, generated protobuf code, etc.). Treat this as the canonical place for cross-service definitions.
- **Kafka/Postgres Flow** — gateway -> chat service -> Postgres persistence + Kafka event -> notification & socket listeners -> WebSocket broadcast.
- **Go Workspace** — `go.work` ties `apps/*` and `pkg/` together for local development.
- **Docker Compose** — split into `docker-compose.infra.yml` (Kafka) and `docker-compose.app.yml` (run the services inside containers).

## Quick Start

1. `docker compose -f docker-compose.infra.yml up -d`
2. In separate shells (or use `docker-compose.app.yml`):
   - `cd apps/chat-service && go run ./cmd/chat-service`
   - `cd apps/notification-service && go run ./cmd/notification-service`
   - `cd apps/socket-service && go run ./cmd/socket-service`
   - `cd apps/gateway && go run ./cmd/gateway`
3. Test the flow:
   ```bash
   curl -X POST http://localhost:8080/chat/messages \
     -H 'Content-Type: application/json' \
     -d '{"senderId":"user-1","receiverId":"user-2","content":"hello"}'
   ```
4. WebSocket endpoint lives at `ws://localhost:9200/ws`.
5. Chat-service relies on manual migrations (GORM auto-migrate was removed). SQL migrations sống ở `apps/chat-service/migrations`.
6. Shortcut commands: root `Makefile` (`make proto`, `make migration-create NAME=...`, `make migration-up`, `make migration-down`) hoặc CLI: `cd apps/chat-service && go run ./cmd/migrate [create|up|down]`.

## Environment

- `.env` at repo root feeds all services (loaded via `pkg/config`) including `CHAT_DATABASE_DSN`.
- For host execution Kafka defaults to `localhost:9094` (mapped from the container) and Postgres to `localhost:5432`. In Docker, services use `kafka:9092` and `gsm-postgres:5432`.

## Directory Walkthrough

```
/cursor           # Meta docs (this file)
/apps             # Service code (gateway, chat-service, notification-service, socket-service)
/pkg              # Shared libraries (config, events, generated protobuf stubs)
/docker-compose.* # Deployment helpers
/go.work          # Go workspace definitions
```
