#!/bin/sh

# Start PostgreSQL (if using local DB)
# pg_ctl start -D /var/lib/postgresql/data

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-clearrouter}; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run database migrations
echo "Running database migrations..."
if [ -f "./db/schema.sql" ]; then
    PGPASSWORD=${DB_PASSWORD:-clearrouter_pass} psql -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-clearrouter} -d ${DB_NAME:-clearrouter} -f ./db/schema.sql
fi

# Run additional migrations
if [ -d "./db/migrations" ]; then
    for migration in ./db/migrations/*.sql; do
        if [ -f "$migration" ]; then
            echo "Running migration: $migration"
            PGPASSWORD=${DB_PASSWORD:-clearrouter_pass} psql -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-clearrouter} -d ${DB_NAME:-clearrouter} -f "$migration"
        fi
    done
fi

# Start nginx in background
nginx -g "daemon off;" &

# Start the Go backend
echo "Starting ClearRouter backend..."
export GIN_MODE=release
export PORT=8080
./main