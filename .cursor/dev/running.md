# Running Services

## Manual Run (Without Docker)

If you prefer running binaries directly on the host, the executables automatically load environment variables from the repository's `.env` file. Simply start each service in separate terminals:

```bash
# Terminal 1: Notification Service
cd apps/notification-service && go run ./cmd/notification-service

# Terminal 2: Chat Service
cd apps/chat-service && go run ./cmd/chat-service

# Terminal 3: Socket Service
cd apps/socket-service && go run ./cmd/socket-service

# Terminal 4: Gateway
cd apps/gateway && go run ./cmd/gateway

# Terminal 5: Auth Service (if needed)
cd apps/auth-service && go run ./cmd/auth-service
```

Override any setting by exporting it before launch:

```bash
export KAFKA_BROKERS=localhost:9094
export CHAT_DATABASE_DSN=postgres://user:pass@localhost:5432/db
```

## Testing the Flow

After starting all services, test the chat creation flow:

```bash
curl -X POST http://localhost:8080/chat/messages \
  -H 'Content-Type: application/json' \
  -d '{"senderId":"user-1","receiverId":"user-2","content":"hello"}'
```

## WebSocket Connection

WebSocket endpoint lives at `ws://localhost:9200/ws`.

Connect using any WebSocket client to receive real-time updates for chat messages and notifications.

## Logging

Logging uses Zerolog with JSON output. Set `LOG_OUTPUT_DIR` (default `./logs`) to control file location. Docker compose sets it to `/var/log/app`, which Promtail watches and ships to Loki/Grafana.

When running manually, logs are written to `./logs/app.log` by default.

## Observing Logs

### Via Grafana (Docker)

1. Start infrastructure: `docker compose -f docker-compose.infra.yml up -d`
2. Access Grafana: `http://localhost:3000` (default `admin/admin`)
3. Add Loki as data source: `http://loki:3100`
4. Query logs using LogQL

### Via Files (Manual Run)

```bash
tail -f logs/app.log
```

### Via Docker Logs

```bash
docker logs -f gsm-gateway
docker logs -f gsm-chat-service
docker logs -f gsm-notification-service
docker logs -f gsm-socket-service
```

