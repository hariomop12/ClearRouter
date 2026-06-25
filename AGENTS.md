# ClearRouter — Agent Guide

## Stack
- **Frontend:** React + TypeScript + Vite (apps/frontend)
- **Backend:** Go + Gin + GORM (apps/backend)
- **Database:** PostgreSQL
- **Auth:** JWT (golang-jwt) + Google/GitHub OAuth
- **Payments:** Razorpay
- **Email:** Resend
- **LLM Providers:** OpenAI, Anthropic, Google/Gemini, DeepSeek, Mistral

## Dev

```sh
# Backend
cd apps/backend && go run ./cmd/server/main.go

# Frontend
cd apps/frontend && pnpm dev

# Docker dev
docker compose -f docker-compose.dev.yml up -d
```

## CI/CD
- `.github/workflows/ci.yml` — Every push: load schema, start backend, Newman smoke tests, build
- `.github/workflows/docker-publish.yml` — Push to main: build Docker image → GHCR → deploy hook

## Newman Tests
```sh
newman run apis/clearrouter.postman_collection.json \
  -e apis/clearrouter.postman_environment.json
```

## Env Vars (all required unless noted)
```
DATABASE_URL=postgres://user:pass@host:5432/clearrouter?sslmode=disable
JWT_SECRET=<random-string>

# OAuth
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=xxx
GOOGLE_REDIRECT_URI=https://host/auth/google/callback
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
GITHUB_REDIRECT_URI=https://host/auth/github/callback
FRONTEND_URL=https://frontend-host

# Payments (Razorpay)
RAZORPAY_KEY_ID=xxx
RAZORPAY_KEY_SECRET=xxx
RAZORPAY_WEBHOOK_SECRET=xxx

# Email (Resend)
RESEND_API_KEY=xxx
RESEND_FROM_EMAIL=noreply@domain

# LLM API Keys (at least one)
OPENAI_API_KEY=xxx
ANTHROPIC_API_KEY=xxx
GOOGLE_API_KEY=xxx          # or GEMINI_API_KEY
DEEPSEEK_API_KEY=xxx
MISTRAL_API_KEY=xxx

# Seed admin (optional, defaults apply)
SEED_DEFAULT_USER_EMAIL=admin@clearrouter.local
SEED_DEFAULT_USER_PASSWORD=admin123
SEED_DEFAULT_USER_NAME=Admin

# Other
GIN_MODE=release
APP_URL=https://host
CURRENCY=INR
USD_TO_INR=83
```

## Build commands
```sh
# Backend
cd apps/backend && go build ./cmd/server
# Frontend
cd apps/frontend && pnpm build
```
