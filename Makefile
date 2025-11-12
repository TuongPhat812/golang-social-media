PROTO_SRC=proto/chat/v1/chat_service.proto
PROTOC?=protoc
CHAT_MIGRATION_DIR=apps/chat-service/migrations
CHAT_DB_DSN?=postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable
MIGRATE_CLI=cd apps/chat-service && LOG_OUTPUT_DIR=$${LOG_OUTPUT_DIR:-logs} CHAT_DATABASE_DSN=$(CHAT_DB_DSN) go run ./cmd/migrate

.PHONY: proto migration-create migration-up migration-down

proto:
	@$(PROTOC) --version >/dev/null
	$(PROTOC) --proto_path=proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(PROTO_SRC)

migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration-create NAME=<snake_case_description>"; exit 1; \
	fi
	$(MIGRATE_CLI) create $(NAME)

migration-up:
	$(MIGRATE_CLI) up

migration-down:
	$(MIGRATE_CLI) down
