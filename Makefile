PROTO_SRC=proto/chat/v1/chat_service.proto
PROTOC?=protoc
MIGRATE?=migrate
CHAT_MIGRATION_DIR=apps/chat-service/migrations
CHAT_DB_DSN?=postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable

.PHONY: proto migration-create migration-up migration-down

proto:
	@$(PROTOC) --version >/dev/null
	$(PROTOC) --proto_path=proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(PROTO_SRC)

migration-create:
	@if [ -z "$(NAME)" ]; then \\
		echo "Usage: make migration-create NAME=<snake_case_description>"; exit 1; \\
	fi
	$(MIGRATE) create -ext sql -dir $(CHAT_MIGRATION_DIR) -seq $(NAME)

migration-up:
	$(MIGRATE) -path $(CHAT_MIGRATION_DIR) -database "$(CHAT_DB_DSN)" up

migration-down:
	$(MIGRATE) -path $(CHAT_MIGRATION_DIR) -database "$(CHAT_DB_DSN)" down
