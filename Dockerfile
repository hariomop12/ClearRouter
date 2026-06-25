FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .
RUN go build -o server ./cmd/server/main.go

# Target for local development
FROM golang:1.25-alpine AS dev
RUN apk add --no-cache git ca-certificates build-base
RUN go install github.com/air-verse/air@latest
RUN go install github.com/amacneil/dbmate@latest
WORKDIR /app
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .
EXPOSE 8080
CMD ["sh", "-c", "dbmate --url \"$DATABASE_URL\" --migrations-dir db/migrations up && air -c .air.toml"]

# Target for production
FROM alpine:latest AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl
ARG DBMATE_VERSION=2.24.0
RUN curl -fsSL https://github.com/amacneil/dbmate/releases/download/v${DBMATE_VERSION}/dbmate-linux-amd64 -o /usr/local/bin/dbmate && chmod +x /usr/local/bin/dbmate
COPY --from=builder /app/server .
COPY apps/backend/db/migrations ./db/migrations
EXPOSE 8080
CMD ["sh", "-c", "dbmate --url \"$DATABASE_URL\" --migrations-dir db/migrations up && exec ./server"]
