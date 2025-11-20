#!/bin/bash
set -e

echo "🛑 Stopping all services..."
docker-compose -f docker-compose.app.yml down
docker-compose -f docker-compose.infra.yml down -v

echo "🧹 Cleaning up volumes..."
docker volume prune -f

echo "🌐 Removing network..."
docker network rm gsm-network 2>/dev/null || echo "Network already removed"

echo "🚀 Starting infrastructure..."
docker-compose -f docker-compose.infra.yml up -d

echo "⏳ Waiting for services to be ready..."
sleep 15

echo "✅ Checking services status..."
docker ps --filter "name=gsm-" --format "table {{.Names}}\t{{.Status}}" | grep -E "postgres|scylla|kafka"

echo "📊 Running PostgreSQL migrations..."
make migration-up

echo "📊 Running ScyllaDB migrations..."
docker exec gsm-scylla-1 cqlsh -f /var/lib/scylla-init/notification_service.cql

echo "✅ Infrastructure reset complete!"
