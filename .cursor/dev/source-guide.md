# How To Read This Codebase

## 1. Start With the Shared Packages (`pkg/`)
- `pkg/config`: understands how environment variables and `.env` are loaded. All services call `config.LoadEnv()` in `cmd`. Know this first to understand runtime settings.
- `pkg/events`: defines the event payloads that move through Kafka. Look at `chat.go` then `notification.go` for the end-to-end data contract.
- `proto/chat/v1` + `pkg/gen/chat/v1`: contain the protobuf schema and generated Go stubs for the gRPC API (gateway ↔ chat-service).

Understanding these shared pieces gives context before diving into any service.

## 2. Follow the Event Flow
For a concrete example, see `.cursor/dev/features/chat/create-chat.md` which follows the end-to-end chat creation flow including database persistence and Kafka side effects.

## 3. Domain Placement
- Each service holds its own `internal/domain/*` package to prevent cross-service coupling.
- Shared representations exist only in `pkg/events` or the generated protobuf packages under `pkg/gen`.

## 4. Configuration & Environment
- `.env` is auto-loaded; inspect it if behaviour looks wrong.
- Services read `KAFKA_BROKERS`, `*_PORT`, consumer group IDs, and now `CHAT_DATABASE_DSN` for Postgres. Check `pkg/config/env.go` for parsing behavior.

## 5. Running / Debugging Tips
- Use the log statements added to Kafka publishers/subscribers to confirm which broker each service uses.
- For quick testing, run a single POST request and observe logs across services.
- If you need to inspect Kafka topics, exec into the Kafka container (`docker exec -it gsm-kafka bash`) và dùng `kafka-console-consumer.sh`.
- Cần xem DB? `docker exec -it gsm-postgres psql -U chat_user -d chat_service`.
- Quản lý migration: `cd apps/chat-service && go run ./cmd/migrate [create|up|down]`.

---
**Need to modify or extend?**
- Add new shared functionality under `pkg/` to keep services slim.
- To add a new service, copy the structure under `apps/<service>` to stay consistent.
- Update `go.work` and `docker-compose.app.yml` when adding or moving services.
