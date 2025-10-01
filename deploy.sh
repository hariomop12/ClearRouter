#!/bin/bash

echo "🚀 ClearRouter Production Deployment Script"
echo "============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Building ClearRouter full-stack application...${NC}"

# Build and start the application
docker-compose -f docker-compose.production.yml down
docker-compose -f docker-compose.production.yml build --no-cache
docker-compose -f docker-compose.production.yml up -d

echo -e "${GREEN}✅ Build complete!${NC}"
echo ""
echo -e "${YELLOW}🔍 Checking application status...${NC}"

# Wait for services to start
sleep 10

# Check if services are running
if docker-compose -f docker-compose.production.yml ps | grep -q "Up"; then
    echo -e "${GREEN}✅ ClearRouter is running successfully!${NC}"
    echo ""
    echo "📱 Application URLs:"
    echo "   Local: http://localhost"
    echo ""
    echo "🔗 To make it accessible via Cloudflare Tunnel:"
    echo "   1. Install cloudflared: https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/"
    echo "   2. Run: cloudflared tunnel --url http://localhost:80"
    echo ""
    echo "📊 To view logs:"
    echo "   docker-compose -f docker-compose.production.yml logs -f"
    echo ""
    echo "🛑 To stop the application:"
    echo "   docker-compose -f docker-compose.production.yml down"
else
    echo -e "${RED}❌ Something went wrong. Check the logs:${NC}"
    docker-compose -f docker-compose.production.yml logs
fi