# Production Dockerfile for React frontend
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY apps/frontend/package*.json apps/frontend/
RUN cd apps/frontend && npm ci --only=production

# Copy source code and build
COPY apps/frontend/ apps/frontend/
RUN cd apps/frontend && npm run build

# Production stage with nginx
FROM nginx:alpine

# Copy built assets
COPY --from=builder /app/apps/frontend/dist /usr/share/nginx/html

# Copy nginx configuration
COPY apps/frontend/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]