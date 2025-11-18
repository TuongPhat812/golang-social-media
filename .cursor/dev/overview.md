# Project At A Glance

**Name:** Golang Social Media Prototype

**Purpose:** Provide a sample social-media backend composed of multiple Go services (gateway, chat, notification, socket) communicating via gRPC and Kafka. The repo demonstrates a DDD-inspired layout inside a mono-repo.

## Core Concepts

- **apps/** — contains all runnable services. Each service follows a `cmd/` entrypoint plus `internal/` layered packages (domain, application, infrastructure, interfaces).
- **pkg/** — reusable building blocks shared across services (`config`, `events`, generated protobuf code, etc.). Treat this as the canonical place for cross-service definitions.
- **Kafka/Postgres Flow** — gateway -> chat service -> Postgres persistence + Kafka event -> notification & socket listeners -> WebSocket broadcast.
- **Go Workspace** — `go.work` ties `apps/*` and `pkg/` together for local development.
- **Docker Compose** — split into `docker-compose.infra.yml` (Kafka, Postgres, Loki, Promtail, Grafana, Kafka UI) and `docker-compose.app.yml` (run the services inside containers).
- **Logging** — Zerolog writes JSON logs to `LOG_OUTPUT_DIR` (default `./logs`); Promtail ships them to Loki/Grafana when running via Docker Compose.

## Quick Start

1. **Setup**: See [setup.md](./setup.md) for development environment setup
2. **Start infrastructure**: `docker compose -f docker-compose.infra.yml up -d` (see [docker.md](./docker.md) for details)
3. **Run services**: See [running.md](./running.md) for manual execution or use `docker-compose.app.yml`
4. **Test the flow**: See [running.md](./running.md) for testing examples
5. **Migrations**: See [migrations.md](./migrations.md) for database migration commands

## Environment

- `.env` at repo root feeds all services (loaded via `pkg/config`). See [environment.md](./environment.md) for all configuration options.
- For host execution Kafka defaults to `localhost:9092` and Postgres to `localhost:5432`. In Docker, services use `kafka:9092` and `gsm-postgres:5432`.

## Directory Walkthrough

```
/cursor           # Meta docs (this file)
/apps             # Service code (gateway, chat-service, notification-service, socket-service)
/pkg              # Shared libraries (config, events, generated protobuf stubs)
/docker-compose.* # Deployment helpers
/go.work          # Go workspace definitions
```
