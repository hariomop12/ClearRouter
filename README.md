# ClearRouter

Full-Stack AI Chat Service built with Go, React, and PostgreSQL in a monorepo architecture.

## Features

### Backend (Go + Gin)
- 🤖 AI Chat API with multiple providers (OpenAI, Google)
- 🔐 JWT Authentication
- 💳 Credit-based usage system
- 🔑 API Key management
- 📊 Chat history tracking

### Frontend (React + TypeScript)
- ⚛️ Modern React with TypeScript
- 🎨 Tailwind CSS for styling
- ⚡ Vite for fast development
- 🔄 Hot module replacement

### Infrastructure
- 🐳 Docker containerization
- 🔄 Database migrations
- 🧪 Test coverage with CI/CD
- 📦 Monorepo structure with workspaces

## Quick Start

### One Command Setup

```bash
# Install all dependencies
npm run install:all

# Start everything in development mode
npm run dev
```

This will start:
- Backend API at `http://localhost:8080`
- Frontend at `http://localhost:3000`
- PostgreSQL database

### Docker Setup

```bash
# Development with hot reload
docker compose -f docker-compose.dev.yml up -d

# Production build
docker-compose up
```

### Database Seeding (dev)

The backend seeds a default user on startup (idempotent) so you can log in immediately during development.

- Default user:
  - Email: `admin@clearrouter.local`
  - Name: `Admin`
  - Password: `admin123`
  - EmailVerified: true

- Environment overrides (set in `docker-compose.dev.yml` under `backend.environment` or your shell):
  - `SEED_DEFAULT_USER_EMAIL` – default `admin@clearrouter.local`
  - `SEED_DEFAULT_USER_NAME` – default `Admin`
  - `SEED_DEFAULT_USER_PASSWORD` – default `admin123`
  - `SEED_ENABLE` – set to `false` to disable seeding

- Apply migrations (if needed):
  - In dev (container): `docker compose -f docker-compose.dev.yml exec backend dbmate -d /app/db/migrations up`
  - From host (psql/TablePlus): use Postgres URI `postgres://clearrouter:clearrouter_pass@localhost:5432/clearrouter?sslmode=disable`

Notes:
- Seeding runs after the backend connects to the DB and only creates the user if it does not already exist.
- Do not use these defaults in production.

### Add Credits (dev)

In development, simulate Razorpay's webhook to credit a user account.

- **Webhook endpoint**: `POST /credits/add`
- **Header**: `X-Razorpay-Signature` (HMAC-SHA256 of raw body using your webhook secret)
- **Helper script**: `scripts/gen_razorpay_signature.sh`

1) Create a payload file (example):
```json
{
  "event": "payment.captured",
  "payload": {
    "payment": { "entity": { "id": "rzp_test_xxx", "amount": 10000, "status": "captured", "order_id": "order_xxx" } },
    "order":   { "entity": { "id": "order_xxx", "amount": 10000, "notes": { "user_id": "<USER_UUID>" } } }
  }
}
```

2) Generate signature and call webhook:
```bash
# Set your webhook secret (from Razorpay dashboard)
export RAZORPAY_WEBHOOK_SECRET="your_webhook_secret"

# Generate signature (prints hex)
SIG=$(scripts/gen_razorpay_signature.sh --payload payload.json)

# Send webhook (use --data-binary to preserve raw bytes)
curl -sS http://localhost:8080/credits/add \
  -H "Content-Type: application/json" \
  -H "X-Razorpay-Signature: $SIG" \
  --data-binary @payload.json
```

3) Verify credits (JWT required):
```bash
# Login to get token (example)
curl -s http://localhost:8080/auth/login -H 'Content-Type: application/json' \
  --data '{"email":"admin@clearrouter.local","password":"admin123"}'

# Then call credits endpoint with Authorization: Bearer <token>
curl -s http://localhost:8080/credits -H "Authorization: Bearer <JWT>"
```

Notes:
- The signature must be computed on the exact raw bytes sent to the server.
- For another user, replace `<USER_UUID>` in the payload `order.entity.notes.user_id` with that user's UUID.
- In production, Razorpay calls this endpoint directly; do not expose your webhook secret.

## Project Structure

```
ClearRouter/
├── apps/
│   ├── backend/                # Go backend application
│   │   ├── cmd/server/        # Application entry point
│   │   ├── internal/          # Private application code
│   │   │   ├── handlers/      # HTTP handlers
│   │   │   ├── middleware/    # HTTP middleware
│   │   │   ├── models/       # Data models
│   │   │   ├── providers/    # AI provider implementations
│   │   │   ├── services/     # Business logic
│   │   │   └── utils/        # Utility functions
│   │   ├── Dockerfile        # Development container
│   │   └── Dockerfile.prod   # Production container
│   └── frontend/              # React frontend application
│       ├── src/              # Source code
│       ├── public/           # Static assets
│       ├── Dockerfile        # Production container
│       ├── Dockerfile.dev    # Development container
│       └── nginx.conf        # Production nginx config
├── db/                       # Database files
│   ├── schema.sql           # Database schema
│   └── migrations/          # Database migrations
├── docker-compose.yml       # Production compose
├── docker-compose.dev.yml   # Development compose
├── package.json             # Root workspace config
└── .env                     # Environment variables
```

## Development Commands

```bash
# Install all dependencies
npm run install:all

# Development (both services)
npm run dev

# Development (individual services)
npm run dev:backend
npm run dev:frontend

# Docker development
npm run dev:docker

# Build
npm run build
npm run build:frontend
npm run build:backend

# Clean all build artifacts
npm run clean
```

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user

### Chat
- `POST /chat` - Send chat message
- `GET /chat/history` - Get chat history

### Models
- `GET /models` - List available AI models

### Credits
- `GET /credits` - Get user credits
- `POST /credits/add` - Add credits (admin)

### API Keys
- `GET /apikeys` - List user's API keys
- `POST /apikeys` - Create new API key
- `DELETE /apikeys/:id` - Delete API key

## Technology Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Authentication**: JWT
- **Migration**: dbmate

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Development**: Hot Module Replacement

### DevOps
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions
- **Development**: Air (Go hot reload)

## Testing

```bash
# Backend tests
cd apps/backend && go test ./...

# Frontend tests
cd apps/frontend && npm test

# Run with coverage
cd apps/backend && go test -cover ./...
```

## Environment Variables

Key environment variables:

```bash
# Backend
DATABASE_URL=postgres://clearrouter:clearrouter_pass@localhost:5432/clearrouter?sslmode=disable
JWT_SECRET=your-secret-key
OPENAI_API_KEY=your-openai-key
GOOGLE_API_KEY=your-google-key
GIN_MODE=debug

# Frontend
VITE_API_URL=http://localhost:8080
```

## Deployment

### Development
```bash
# Start with hot reload
npm run dev

# Or with Docker
docker-compose -f docker-compose.dev.yml up
```

### Production
```bash
# Build and start
docker-compose up -d

# Manual migrations (if needed)
docker-compose exec backend dbmate up
```

### CI/CD Pipeline

GitHub Actions pipeline includes:
- Go testing and linting
- React build and testing
- Docker image building
- Security scanning

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details.