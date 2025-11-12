# Notes for Future AI Assistants

1. **Always call `config.LoadEnv()`**: Every service already does this in `cmd/main.go`. If you create new binaries, follow the same pattern so `.env` and env overrides work.
2. **Respect Service Boundaries**:
   - Domain structs live under `apps/<service>/internal/domain`. Do not import one service’s domain package into another; use `pkg/events` or protobuf stubs under `pkg/gen` for shared shapes.
   - When adding cross-service contracts, define/update `.proto` files under `proto/` and regenerate Go code into `pkg/gen`.
3. **Kafka Considerations**:
   - Kafka is expected to run via `docker-compose.infra.yml` (single-node KRaft). Broker hosts: `localhost:9094` (host), `kafka:9092` (containers).
   - Prefer `github.com/segmentio/kafka-go` APIs already in use. Ensure replication factors stay at 1 unless you expand infra.
4. **gRPC Contracts**:
   - Use the protobuf definitions in `proto/` and generated Go code under `pkg/gen`. Regenerate with `protoc` (see tooling in `.tools/`) when the schema changes.
5. **Go Modules & Workspace**:
   - Keep modules under `apps/<service>` and `pkg`. If you add new modules, update `go.work` and add the corresponding `replace` statements (`pkg` should point via `../../pkg`). Run `go work sync` afterwards.
6. **Repository Hygiene**:
   - `.gitignore` already excludes `.env` and Docker artifacts. Keep sensitive info out of the repo.
   - When restructuring, run `gofmt` and `go mod tidy` per module.
7. **Documentation**:
   - Update `README.md` if you add services or change the event flow.
   - Add notes to `.cursor/` when making significant structural decisions.

Following these guardrails will keep the codebase clean and consistent with the current architecture.
