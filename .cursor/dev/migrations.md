# Database Migrations

## PostgreSQL (Chat Service)

Chat-service connects to Postgres using the `CHAT_DATABASE_DSN` environment variable (defaults to `postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable`).

Schema changes rely on explicit SQL migrations (GORM auto-migrate is disabled).

Versioned SQL migrations live in `apps/chat-service/migrations`.

### Using Makefile (from repo root)

```bash
# create a new migration (NAME is required)
make migration-create NAME=add_new_table

# apply migrations (uses CHAT_DB_DSN or defaults to local postgres)
make migration-up

# rollback one step
make migration-down
```

### Using Go CLI (from chat-service directory)

```bash
cd apps/chat-service
go run ./cmd/migrate create add_messages_index   # generates stub files
go run ./cmd/migrate up                          # applies pending migrations
go run ./cmd/migrate down                        # rolls back one step
```

Then edit the generated `.up.sql` / `.down.sql` files.

## ScyllaDB (Notification Service)

Schema changes live in `apps/notification-service/infra/scylla/notification_service.cql`.

### Apply Migrations Locally

If you have `cqlsh` installed on your host machine:

```bash
cqlsh -f apps/notification-service/infra/scylla/notification_service.cql localhost 9042
```

The command assumes a local Scylla node on port `9042` (matching `docker-compose.infra.yml`). For remote clusters, adjust the host/port and authentication flags as needed.

**Note**: If you installed `cqlsh` via `snap`, you may encounter permission issues reading files from the project directory. Use the Docker exec method below instead.

### Using Docker exec (recommended for snap-installed cqlsh)

If you installed `cqlsh` via `snap`, it may have permission issues reading files from the project directory. Use Docker exec instead:

```bash
# Initial schema setup
docker exec -it gsm-scylla cqlsh -f /var/lib/scylla-init/notification_service.cql

# Apply migrations (e.g., add read_at column)
docker exec -it gsm-scylla cqlsh -f /app/infra/scylla/add_read_at_column.cql
```

**Notes**:
- When running `cqlsh` from inside the container, you don't need to specify `localhost 9042` - it will use the default connection automatically.
- The keyspace uses `replication_factor=3` (production best practice). With a single-node setup, ScyllaDB will still create the keyspace but data will only be stored on the single node. For full replication benefits, set up a 3-node cluster.
- If you encounter "Unknown identifier" errors, you may need to apply migrations. See migration files in `apps/notification-service/infra/scylla/`.

### ScyllaDB Web UI (Cassandra Web UI)

After starting the infrastructure stack, access the Cassandra Web UI at `http://localhost:8083` to:

- Browse keyspaces and tables
- Execute CQL queries interactively
- View and edit table data
- Export/import data
- Explore the database schema

The UI automatically connects to the ScyllaDB instance running in Docker. No authentication is required for the local development setup.

