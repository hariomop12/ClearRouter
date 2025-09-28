# Development stage
FROM golang:1.23-alpine

# Install development dependencies
RUN apk add --no-cache git curl gcc musl-dev postgresql-client

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@v1.42.0

# Install dbmate for database migrations
RUN curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 && \
    chmod +x /usr/local/bin/dbmate

WORKDIR /app

# Copy go modules files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project structure
COPY . .

# Make sure the entrypoint script is executable
RUN chmod +x /app/apps/backend/docker-entrypoint.sh

# Use the entrypoint script
ENTRYPOINT ["/bin/sh", "/app/apps/backend/docker-entrypoint.sh"]
CMD ["go", "run", "/app/apps/backend/cmd/server/main.go"]