PROTO_SRC=proto/chat/v1/chat_service.proto
PROTOC?=protoc
CHAT_MIGRATION_DIR=apps/chat-service/migrations
CHAT_DB_DSN?=postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable
MIGRATE_CLI=cd apps/chat-service && LOG_OUTPUT_DIR=$${LOG_OUTPUT_DIR:-logs} CHAT_DATABASE_DSN=$(CHAT_DB_DSN) go run ./cmd/migrate

.PHONY: proto migration-create migration-up migration-down setup-ubuntu-deps load-test-chat

proto:
	@$(PROTOC) --version >/dev/null
	@mkdir -p pkg/gen
	@export PATH=$$(go env GOPATH)/bin:$$PATH && $(PROTOC) --proto_path=proto --go_out=paths=source_relative:pkg/gen --go-grpc_out=paths=source_relative:pkg/gen $(PROTO_SRC)

migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration-create NAME=<snake_case_description>"; exit 1; \
	fi
	$(MIGRATE_CLI) create $(NAME)

migration-up:
	$(MIGRATE_CLI) up

migration-down:
	$(MIGRATE_CLI) down

setup-ubuntu-deps:
	sudo apt update
	sudo apt install -y snapd protobuf-compiler
	sudo snap install cqlsh
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

scylla-migrate-add-read-at:
	@echo "Applying ScyllaDB migration: add_read_at_column..."
	@docker exec -i gsm-scylla-1 cqlsh -e "USE notification_service; ALTER TABLE notifications_by_user ADD read_at timestamp;" || echo "Note: If column already exists, this error is expected."
	@echo "Migration completed!"

load-test-chat:
	@echo "Running load test on chat endpoint..."
	@go run scripts/load_test_chat.go
