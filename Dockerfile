# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY apps/frontend/package*.json ./
RUN npm ci
COPY apps/frontend .
RUN npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.* ./
COPY apps/backend ./apps/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./apps/backend/cmd/server/main.go

# Final image
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy backend binary
COPY --from=backend-builder /app/server .
COPY apps/backend/db/schema.sql ./db/schema.sql
COPY apps/backend/db/migrations ./db/migrations

# Copy nginx configuration
COPY apps/frontend/nginx.conf /etc/nginx/nginx.conf

# Install nginx
RUN apk add --no-cache nginx

EXPOSE 80 8080

# Start script
COPY --chmod=755 <<'EOF' /start.sh
#!/bin/sh
nginx &
./server
EOF

CMD ["/start.sh"]