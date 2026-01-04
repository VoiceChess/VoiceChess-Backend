#!/bin/bash

# VoiceChess Backend Deployment Script
# Usage: ./deploy.sh

set -e

echo "🚀 Starting VoiceChess Backend Deployment..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if .env.docker exists
if [ ! -f .env.docker ]; then
    echo -e "${RED}❌ Error: .env.docker file not found!${NC}"
    echo "Please create .env.docker with your production environment variables"
    exit 1
fi

# Load environment variables from .env.docker
export $(cat .env.docker | grep -v '^#' | xargs)

echo -e "${YELLOW}📦 Stopping existing containers...${NC}"
docker-compose down

echo -e "${YELLOW}🏗️  Building Docker image...${NC}"
docker-compose build --no-cache

echo -e "${YELLOW}🚀 Starting containers...${NC}"
docker-compose up -d

echo -e "${GREEN}✅ Deployment completed!${NC}"
echo -e "${GREEN}🔍 Checking container status...${NC}"
docker-compose ps

echo -e "${GREEN}📋 Viewing logs (press Ctrl+C to exit)...${NC}"
docker-compose logs -f
