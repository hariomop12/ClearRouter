#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until pg_isready -h postgres -U clearrouter; do
    echo "Database not ready yet... waiting"
    sleep 2
done

echo "Database is ready!"

# Set default DATABASE_URL if not provided
if [ -z "$DATABASE_URL" ]; then
    export DATABASE_URL="postgres://clearrouter:clearrouter_pass@postgres:5432/clearrouter?sslmode=disable"
fi

# Migrations are handled manually by the user
# To run migrations: docker compose exec backend dbmate up
echo "Skipping automatic migrations. Run manually with: docker compose exec backend dbmate up"

echo "Starting the application..."
# Check if we're in development mode (with full source tree) or production mode (with just binary)
if [ -d "/app/apps/backend" ]; then
    cd /app/apps/backend
else
    cd /app
fi
exec "$@"