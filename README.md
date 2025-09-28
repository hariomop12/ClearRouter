# ClearRouter - AI Chat Service 🚀

ClearRouter is a robust, scalable AI chat service built with Go, featuring multiple AI provider integrations, user authentication, credit system, and comprehensive chat history management.

## 🏗️ Architecture Overview

### Tech Stack
- **Backend**: Go 1.23
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 15 with GORM ORM
- **Containerization**: Docker + Docker Compose
- **Hot Reloading**: Air (development)
- **Database Migrations**: dbmate
- **Authentication**: JWT + API Key based
- **Payment Integration**: Razorpay
- **AI Providers**: OpenAI, Google

### Project Structure

```
ClearRouter/
├── apps/
│   └── backend/
│       ├── cmd/
│       │   └── server/          # Application entry point (main.go)
│       ├── internal/
│       │   ├── handlers/        # HTTP request handlers
│       │   │   ├── auth.go      # Authentication endpoints
│       │   │   ├── apikey.go    # API key management
│       │   │   ├── chat.go      # Chat completions
│       │   │   ├── chat_history.go # Chat history management
│       │   │   ├── credits.go   # Credit system & payments
│       │   │   └── models.go    # Model listing
│       │   ├── middleware/      # HTTP middleware
│       │   ├── models/          # Database models & DTOs
│       │   │   ├── user.go      # User model
│       │   │   ├── apikey.go    # API key model
│       │   │   ├── credits.go   # Credits model
│       │   │   ├── chat_history.go # Chat & message models
│       │   │   └── *_models.go  # Provider-specific models
│       │   ├── providers/       # AI provider integrations
│       │   │   ├── openai_provider.go
│       │   │   └── google_provider.go
│       │   ├── services/        # Business logic
│       │   └── utils/           # Utility functions (JWT, email)
│       ├── .air.toml           # Hot reload configuration
│       ├── Dockerfile          # Docker build instructions
│       └── docker-entrypoint.sh # Container startup script
├── db/
│   ├── schema.sql              # Database schema (auto-generated)
│   └── migrations/             # Database migrations
│       ├── 20250927072843_m1.sql
│       ├── 20250927101618_m2.sql
│       ├── 20250927102719_fix_api_keys_schema.sql
│       ├── 20250927192900_add_unique_constraint_to_credits.sql
│       └── 20250928134500_create_chat_history.sql
├── docker-compose.yml          # Production Docker configuration
├── docker-compose.dev.yml      # Development Docker configuration
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
└── .env.example               # Environment variables template
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### Setup & Run
```bash
# Clone the repository
git clone https://github.com/hariomop12/ClearRouter.git
cd ClearRouter

# Copy environment variables
cp .env.example .env

# Start everything (development mode with hot reloading)
docker compose -f docker-compose.dev.yml up --build

# Or for production
docker compose up --build
```

**That's it! 🎉** The entire application will start:
- **PostgreSQL Database**: Running on port 5432
- **Go Backend API**: Running on port 8080
- **Database Migrations**: Applied automatically
- **Hot Reloading**: Enabled in development mode

### Database Connection
```
Host: localhost
Port: 5432
Database: clearrouter
Username: clearrouter
Password: clearrouter_pass
URL: postgres://clearrouter:clearrouter_pass@localhost:5432/clearrouter
```

## 🏗️ Core Components

### 1. Authentication System
- **JWT Authentication**: For user sessions and protected routes
- **API Key Authentication**: For chat completions API (OpenAI-compatible)
- **Email Verification**: Secure user registration
- **Bcrypt Password Hashing**: Industry-standard security

### 2. Database Architecture
- **User Management**: Account creation, verification, authentication
- **API Key System**: Generate and manage API keys for services
- **Credit System**: Pay-per-use model with Razorpay integration
- **Chat History**: Persistent conversation storage
- **Usage Tracking**: API usage logs and analytics

### 3. AI Provider System
- **Multi-Provider Support**: OpenAI, Google AI
- **Unified Interface**: Single API for multiple providers
- **Automatic Routing**: Load balancing and failover
- **Cost Optimization**: Smart provider selection

## 📡 API Endpoints

### 🔐 Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/auth/signup` | User registration | ❌ |
| `GET` | `/auth/verify` | Email verification | ❌ |
| `POST` | `/auth/login` | User login | ❌ |

### 🔑 API Key Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/keys/create` | Create new API key | ✅ JWT |
| `GET` | `/keys` | List user's API keys | ✅ JWT |

### 💳 Credits & Payments
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/credits/order` | Create credit purchase order | ✅ JWT |
| `POST` | `/credits/add` | Razorpay payment webhook | ❌ |
| `GET` | `/credits` | Get user's credit balance | ✅ JWT |

### 🤖 Chat Completions (OpenAI Compatible)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/v1/chat/completions` | AI chat completions | ✅ API Key |
| `GET` | `/models` | List available AI models | ❌ |

### 📝 Chat History Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/newchat` | Create new chat session | ✅ JWT |
| `GET` | `/chathistory` | Get user's chat history | ✅ JWT |
| `GET` | `/chathistory/:chatId` | Get specific chat details | ✅ JWT |
| `DELETE` | `/chathistory/:chatId` | Delete a chat session | ✅ JWT |

### 📊 General
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/` | Health check | ❌ |

## 🔧 Environment Variables

Create a `.env` file from `.env.example`:

```bash
# Database Configuration
DATABASE_URL=postgres://clearrouter:clearrouter_pass@postgres:5432/clearrouter?sslmode=disable

# Application Configuration
GIN_MODE=debug
GORM_DEBUG=true

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here

# AI Provider API Keys
OPENAI_API_KEY=your_openai_api_key_here
GOOGLE_API_KEY=your_google_api_key_here

# Payment Configuration (Razorpay)
RAZORPAY_KEY_ID=your_razorpay_key_id
RAZORPAY_KEY_SECRET=your_razorpay_key_secret

# Email Configuration (optional)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_email@example.com
SMTP_PASSWORD=your_email_password
```

#### Provider Service
- Manages multiple AI provider integrations
- Supports OpenAI and Google AI
- Extensible architecture for adding new providers

#### Credits Service
- Manages user credits
- Handles credit deduction for API usage
- Integrates with payment system

#### Chat Service
- Manages chat sessions
- Handles message routing to appropriate AI providers
- Maintains chat history

## Security Features

1. **Authentication**
   - JWT token-based authentication
   - Secure password hashing
   - Email verification

2. **API Security**
   - API key authentication
   - Rate limiting
   - Usage tracking

## 💻 Usage Examples

### 1. User Registration & Login
```bash
# Register a new user
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"securepass123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"securepass123"}'
```

### 2. Create API Key
```bash
curl -X POST http://localhost:8080/keys/create \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### 3. Chat Completions (OpenAI Compatible)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

## 🚀 Development

### Local Development
```bash
# Start in development mode (with hot reloading)
docker compose -f docker-compose.dev.yml up --build

# View logs
docker compose -f docker-compose.dev.yml logs -f backend

# Stop services
docker compose -f docker-compose.dev.yml down -v
```

### Production Deployment
```bash
# Start in production mode
docker compose up --build -d

# View logs
docker compose logs -f backend

# Stop services
docker compose down
```

## 🗄️ Database Schema

The system uses PostgreSQL with automatic migrations:

- **users**: User accounts with UUID-based identifiers
- **api_keys**: API key management with usage tracking
- **credits**: Credit balance and transaction history
- **payments**: Payment processing records
- **chats**: Chat session management
- **chat_history_messages**: Individual chat messages
- **api_usage_logs**: API usage analytics

## 🛠️ Tech Features

- **🔥 Hot Reloading**: Air for instant development feedback
- **🐳 Containerized**: Full Docker support
- **📊 Auto Migrations**: Database schema managed automatically
- **🔐 Security**: JWT + API Key authentication
- **💳 Payments**: Razorpay integration
- **📝 Logging**: Comprehensive request/response logging
- **🚀 Scalable**: Microservice-ready architecture

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Hariom** - [@hariomop12](https://github.com/hariomop12)

---

**🎉 Happy Coding!** If you found this project helpful, please give it a ⭐️
   - User association

4. **chat_history**
   - Chat session management
   - Message storage
   - User association

## Development Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Run database migrations:
   ```bash
   dbmate up
   ```
4. Start the development server:
   ```bash
   docker-compose -f docker-compose.dev.yml up
   ```

## Deployment

The application is containerized using Docker and can be deployed using the provided docker-compose configurations:

- `docker-compose.yml`: Production configuration
- `docker-compose.dev.yml`: Development configuration

## Future Considerations

1. **Scaling**
   - Horizontal scaling of API servers
   - Database replication
   - Caching layer implementation

2. **Monitoring**
   - API usage metrics
   - Error tracking
   - Performance monitoring

3. **Features**
   - Additional AI providers
   - Enhanced chat features
   - Advanced analytics

## 🔄 CI/CD Pipeline

The project includes a simple GitHub Actions pipeline that:

### **Automated Testing**
- ✅ Runs on every push and pull request
- ✅ Sets up PostgreSQL test database
- ✅ Runs Go tests with coverage
- ✅ Builds the application

### **Docker Integration**
- ✅ Builds production Docker image
- ✅ Pushes to Docker Hub (on main branch)
- ✅ Tags with commit SHA and 'latest'

### **Setup Instructions**
1. Add these secrets to your GitHub repository:
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_PASSWORD`: Your Docker Hub password

2. The pipeline will automatically:
   - Test your code on every commit
   - Build and push Docker images on main branch
   - Provide deployment-ready artifacts

### **Pipeline Status**
[![CI/CD](https://github.com/hariomop12/ClearRouter/workflows/ClearRouter%20CI/CD/badge.svg)](https://github.com/hariomop12/ClearRouter/actions)