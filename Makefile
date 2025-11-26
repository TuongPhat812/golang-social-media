PROTO_SRC=proto/chat/v1/chat_service.proto
PROTOC?=protoc
CHAT_MIGRATION_DIR=apps/chat-service/migrations
CHAT_DB_DSN?=postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable
MIGRATE_CLI=cd apps/chat-service && LOG_OUTPUT_DIR=$${LOG_OUTPUT_DIR:-logs} CHAT_DATABASE_DSN=$(CHAT_DB_DSN) go run ./cmd/migrate

.PHONY: proto migration-create migration-up migration-down setup-ubuntu-deps load-test-chat test test-auth test-chat test-gateway test-cover

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

test:
	@echo "Running all tests..."
	@cd apps/auth-service && go test -v ./internal/domain/... ./internal/infrastructure/jwt/test/... ./internal/application/command/test/... ./internal/application/query/test/... || true
	@cd apps/chat-service && go test -v ./... || true
	@cd apps/gateway && go test -v ./... || true
	@cd pkg && go test -v ./... || true

test-auth:
	@echo "Running auth-service tests..."
	@cd apps/auth-service && go test -v ./internal/domain/... ./internal/infrastructure/jwt/test/... ./internal/application/command/test/... ./internal/application/query/test/... || true

test-chat:
	@echo "Running chat-service tests..."
	@cd apps/chat-service && go test -v ./... || true

test-gateway:
	@echo "Running gateway tests..."
	@cd apps/gateway && go test -v ./... || true

test-cover:
	@echo "Running tests with coverage..."
	@cd apps/auth-service && go test -cover ./internal/domain/... ./internal/infrastructure/jwt/test/... ./internal/application/command/test/... ./internal/application/query/test/... || true
	@cd apps/chat-service && go test -cover ./... || true
	@cd apps/gateway && go test -cover ./... || true
	@cd pkg && go test -cover ./... || true

test-domain:
	@echo "Running domain tests only..."
	@cd apps/auth-service && go test -v ./internal/domain/user/test ./internal/domain/role/test ./internal/domain/permission/test ./internal/domain/user_role/test ./internal/domain/role_permission/test ./internal/domain/factories/test || true

test-jwt:
	@echo "Running JWT service tests..."
	@cd apps/auth-service && go test -v ./internal/infrastructure/jwt/test || true

test-command:
	@echo "Running command tests..."
	@cd apps/auth-service && go test -v ./internal/application/command/test || true

test-query:
	@echo "Running query tests..."
	@cd apps/auth-service && go test -v ./internal/application/query/test || true
