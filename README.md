# Golang Social Media Prototype

A microservices-based social media backend built with Go, demonstrating Domain-Driven Design (DDD), event-driven architecture, and modern observability practices.

## Overview

This repository is organised as a mono-repo with multiple Go services under `apps/` and shared packages in `pkg/`. Each service follows a DDD-inspired layout with clear separation between domain, application, infrastructure, and interface layers.

## Services

- **`apps/auth-service`**: HTTP service for user registration, login, profile management, and publishing user events
- **`apps/gateway`**: Gin HTTP gateway orchestrating downstream service calls
- **`apps/chat-service`**: gRPC service for creating chat messages, persisting to Postgres, and emitting Kafka events
- **`apps/notification-service`**: Event consumer service that replicates users to ScyllaDB, creates notifications, and emits follow-up events
- **`apps/socket-service`**: WebSocket service broadcasting events to connected clients

## Project Structure

```
.
├── apps/              # Service implementations
│   ├── auth-service/
│   ├── gateway/
│   ├── chat-service/
│   ├── notification-service/
│   └── socket-service/
├── pkg/               # Shared packages (config, events, protobuf stubs)
├── proto/             # Protobuf schemas
└── .cursor/dev/       # Detailed documentation
```

Each service follows this structure:

```
<service>/
├── cmd/<service-name>/main.go     # Composition root / entrypoint
└── internal/
    ├── application/               # Application services and use-cases
    ├── domain/                    # Domain entities and business logic
    ├── infrastructure/            # Framework adapters (gRPC, Kafka, DB, etc.)
    └── interfaces/                # Transport layers (HTTP, gRPC, WebSocket)
```

## Documentation

Detailed guides are available in `.cursor/dev/`:

- **[Setup Guide](.cursor/dev/setup.md)** - Development environment setup, protobuf generation, Go workspace
- **[Docker Guide](.cursor/dev/docker.md)** - Docker Compose setup, infrastructure services, application containers
- **[Migrations Guide](.cursor/dev/migrations.md)** - Database migrations for PostgreSQL and ScyllaDB
- **[Environment Variables](.cursor/dev/environment.md)** - Configuration and environment variable reference
- **[Running Services](.cursor/dev/running.md)** - Manual execution, testing, logging, and debugging
- **[Event Flow](.cursor/dev/event-flow.md)** - Event-driven architecture, Kafka topics, consumer groups
- **[Source Guide](.cursor/dev/source-guide.md)** - How to read and navigate the codebase
- **[Overview](.cursor/dev/overview.md)** - Project at a glance and quick start

## Quick Start

1. **Setup development environment**: See [setup.md](.cursor/dev/setup.md)
2. **Start infrastructure**: See [docker.md](.cursor/dev/docker.md)
3. **Run services**: See [running.md](.cursor/dev/running.md)

## Key Technologies

- **Go 1.22+** with Go workspace (`go.work`)
- **gRPC** for inter-service communication
- **Kafka** for event streaming
- **PostgreSQL** for relational data (chat-service)
- **ScyllaDB** for NoSQL data (notification-service)
- **Docker Compose** for local development
- **Grafana + Loki** for log aggregation
- **Prometheus** for metrics
- **Cassandra Web UI** for ScyllaDB management

## Makefile Commands

- `make proto` - Generate protobuf Go code
- `make migration-create NAME=<name>` - Create a new database migration
- `make migration-up` - Apply pending migrations
- `make migration-down` - Rollback one migration
- `make setup-ubuntu-deps` - Install Ubuntu dependencies (protoc, cqlsh, etc.)

## License

This is a learning project. Use it as a reference for building microservices with Go.
