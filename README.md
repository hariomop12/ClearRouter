# ClearRouter

AI chat gateway — proxy 66+ models from OpenAI, Google, Anthropic, Mistral, and DeepSeek through a single API. Credit-based metering, API key management, analytics, and a dashboard.

## Stack

| Layer     | Tech                              |
| --------- | --------------------------------- |
| Backend   | Go + Gin + GORM + PostgreSQL      |
| Frontend  | React 19 + TypeScript + Vite + Tailwind |
| Payments  | Razorpay                          |
| Auth      | JWT + API Keys                    |
| Email     | Resend                            |

## Quick Start

### Prerequisites

- Go 1.24+
- pnpm 8+
- Docker + Docker Compose
- PostgreSQL 15 (or use a remote DB)

### Dev Setup

```bash
# 1. Install frontend deps
pnpm install

# 2. Copy env and configure
cp .env .env.local   # edit as needed

# 3. Start backend (Docker with hot reload via Air)
docker compose -f docker-compose.dev.yml up

# 4. Start frontend (separate terminal)
pnpm dev
```

- Backend API: `http://localhost:8080`
- Frontend: `http://localhost:5173`

Frontend proxies `/api`, `/auth`, `/v1`, etc. to the backend via Vite.

### Default Dev User

| Field    | Value                      |
| -------- | -------------------------- |
| Email    | admin@clearrouter.local    |
| Password | admin123                   |

Seeded automatically on backend startup.

### Production Build

```bash
docker build -t clearrouter-backend .
docker run -p 8080:8080 --env-file .env clearrouter-backend
```

Frontend is deployed separately via Vercel (see `vercel.json`).

## Project Structure

```
ClearRouter/
├── apps/
│   ├── backend/              # Go API server
│   │   ├── cmd/server/       # Entry point
│   │   ├── internal/         # Handlers, models, providers, services, utils
│   │   ├── db/               # SQL migrations
│   │   └── Dockerfile.dev    # Dev container with Air
│   └── frontend/             # React SPA
│       ├── src/
│       └── nginx.conf
├── db/                       # DB schema + migrations
├── scripts/                  # Razorpay webhook testing
├── docker-compose.dev.yml    # Backend dev
├── docker-compose.yml        # Production compose
├── Dockerfile                # Production backend image
└── .env                      # Environment variables
```

## API Overview

| Route                     | Auth     | Description            |
| ------------------------- | -------- | ---------------------- |
| `POST /auth/signup`       | Public   | Register               |
| `POST /auth/login`        | Public   | Login                  |
| `POST /v1/chat/completions` | API Key | Chat completions       |
| `POST /chat`              | JWT      | Dashboard chat         |
| `GET /models`             | Public   | List available models  |
| `GET /credits`            | JWT      | Credit balance         |
| `POST /credits/order`     | JWT      | Create Razorpay order  |
| `GET /analytics/usage`    | JWT      | Usage stats            |
| `POST /keys/create`       | JWT      | Create API key         |
| `GET /health`             | Public   | Health check           |
