# Docker Compose Setup

Two Compose files are provided:

- `docker-compose.infra.yml` — brings up Kafka, Postgres, Loki, Promtail, Kafka UI, Grafana, Prometheus, ScyllaDB (3-node cluster), and Cassandra Web UI.
- `docker-compose.app.yml` — runs the Go services inside containers (expects the infra stack to be running).

## Start Infrastructure

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.infra.yml up -d
```

This brings up Postgres, Kafka, and the observability stack on the shared `gsm-network`:

- **Postgres**: `localhost:5432` (host) / `gsm-postgres:5432` (containers)
- **Kafka**: `localhost:9092` (host) / `kafka:9092` (containers)
- **ScyllaDB Cluster** (3-node):
  - Node 1: `localhost:9042` (host) / `scylla-1:9042` (containers)
  - Node 2: `localhost:9043` (host) / `scylla-2:9042` (containers)
  - Node 3: `localhost:9044` (host) / `scylla-3:9042` (containers)
- **Loki**: `http://localhost:3100`
- **Promtail**: scrapes `/var/log/app/*.log` and forwards to Loki
- **Prometheus**: `http://localhost:9090` (scrapes ScyllaDB metrics)
- **Grafana**: `http://localhost:3000` (default `admin/admin`)
- **Kafka UI**: `http://localhost:8088`
- **Cassandra Web UI**: `http://localhost:8083` (query and browse ScyllaDB data)

## Start Application Services

```bash
cd /home/ubuntu/Workspace/myself/golang-social-media
docker compose -f docker-compose.app.yml up --build
```

Each service runs with code mounted from the host, using Go's `GOTOOLCHAIN=auto` to pull the correct toolchain. The containers connect to Kafka/Postgres and write logs to `/var/log/app` (mounted from `./logs`) so Promtail/Loki can ingest them.

## Stop Services

```bash
docker compose -f docker-compose.app.yml down
docker compose -f docker-compose.infra.yml down
```

## Network Configuration

When running outside Docker, use the host listener `localhost:9092` for Kafka and `localhost:5432` for Postgres. Services inside Docker should use `kafka:9092` and `gsm-postgres:5432`.

For ScyllaDB, when running services manually (outside Docker), connect to all 3 nodes:
- `SCYLLA_HOSTS=localhost:9042,localhost:9043,localhost:9044`

## ScyllaDB Cluster Status

To check the cluster status after starting:

```bash
docker exec -it gsm-scylla-1 nodetool status
```

This will show all 3 nodes and their status (UN = Up Normal).

