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
docker-compose -f docker-compose.dev.yml up

# Production build
docker-compose up
```

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