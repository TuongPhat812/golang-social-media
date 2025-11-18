# Development Setup

## Bootstrap Tooling (Ubuntu)

Run `make setup-ubuntu-deps` once on a fresh machine to install OS-level tools the project relies on:

```bash
make setup-ubuntu-deps
```

This installs:
- `snapd` and `protobuf-compiler` via `apt`
- `cqlsh` via `snap`
- Go protobuf/grpc plugins via `go install`

Make sure `$GOBIN` (or `$GOPATH/bin`) is on your `PATH` so the generators are discoverable.

## Protobuf Generation

Schemas live under `proto/`. Currently only `proto/chat/v1/chat_service.proto`.

Generated Go code is committed under `pkg/gen/chat/v1`.

Regenerate Go code with:

```bash
make proto
```

This requires `protoc` plus the Go plugins (`protoc-gen-go`, `protoc-gen-go-grpc`) on your `PATH`.

## Go Workspace

The repository uses a Go workspace (`go.work`) that includes every service under `apps/` plus the shared `pkg/` module.

Shared code is pulled from `pkg/` via local replace directives, so nothing needs to be published externally.

