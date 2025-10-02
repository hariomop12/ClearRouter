#!/bin/bash

# ClearRouter Database Migration Helper
# This script runs dbmate commands with the correct DATABASE_URL for localhost

# Set the correct DATABASE_URL for host machine (not Docker internal network)
export DATABASE_URL="postgres://clearrouter:clearrouter_pass@localhost:5432/clearrouter?sslmode=disable"

# Check if containers are running
if ! docker compose -f docker-compose.dev.yml ps | grep -q "clearrouter-postgres-1.*Up"; then
    echo "❌ PostgreSQL container is not running!"
    echo "🚀 Start it with: docker compose -f docker-compose.dev.yml up -d"
    exit 1
fi

echo "🔄 Running dbmate command: $*"
dbmate "$@"