FROM golang:1.25-alpine AS builder

WORKDIR /app
# Install Air inside the development stage
RUN go install github.com/air-verse/air@latest

COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .

# Target for local development
FROM builder AS dev
CMD ["air", "-c", ".air.toml"]

# Target for production
FROM alpine:latest AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
