#!/bin/bash

# Telegram Bot Unified Startup Script
# This script handles everything: environment setup, dependencies, and bot startup

export TG_TOKEN=${TG_TOKEN:-TG_TOKEN}
export OPENAI_API_KEY=${OPENAI_API_KEY:-OPENAI_API_KEY}

echo "🤖 Starting Telegram Bot..."

# Check if required files exist
if [ ! -f "docker-compose.yaml" ]; then
    echo "❌ Error: docker-compose.yaml not found in current directory"
    exit 1
fi

# Load .env file if it exists
if [ -f ".env" ]; then
    echo "📂 Loading environment variables from .env file..."
    export $(grep -v '^#' .env | grep -v '^$' | xargs)
fi

# Database configuration
export POSTGRES_USER=${POSTGRES_USER:-postgres}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
# Use Docker service name when running in container
export DATABASE_URL=${DATABASE_URL:-postgresql://postgres:postgres@postgres:5432/postgres}

# Redis configuration
# Use Docker service name when running in container
export REDIS_URL=${REDIS_URL:-redis://redis:6379}

# MinIO/S3 Storage configuration
# Use Docker service name when running in container
export MINIO_ENDPOINT=${MINIO_ENDPOINT:-minio:9000}
export MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-minioadmin}
export MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-minioadmin}
export MINIO_USE_SSL=${MINIO_USE_SSL:-false}
export MINIO_BUCKET=${MINIO_BUCKET:-telegram-images}
export MINIO_PUBLIC_DOMAIN=${MINIO_PUBLIC_DOMAIN:-localhost:9000}

# Telegramify configuration
# Use Docker service name when running in container
export TELEGRAMIFY_URL=${TELEGRAMIFY_URL:-http://telegramify:8000}

echo ""
echo "🔧 Environment Configuration:"
echo "  Database URL: $DATABASE_URL"
echo "  Redis URL: $REDIS_URL"
echo "  MinIO Endpoint: $MINIO_ENDPOINT"
echo "  MinIO Bucket: $MINIO_BUCKET"
echo "  MinIO Public Domain: $MINIO_PUBLIC_DOMAIN"
echo "  Telegramify URL: $TELEGRAMIFY_URL"
echo "  TG Token: ${TG_TOKEN:0:10}***"
echo "  OpenAI API Key: ${OPENAI_API_KEY:0:10}***"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

echo ""
echo "🐳 Setting up Docker environment..."

# Create dokploy-network if it doesn't exist
if ! docker network inspect dokploy-network > /dev/null 2>&1; then
    echo "🌐 Creating dokploy-network..."
    docker network create dokploy-network
fi

# Start dependencies first
echo "📦 Starting dependencies (PostgreSQL, Redis, MinIO, Telegramify)..."
docker-compose up -d postgres redis minio telegramify

# Wait for services to be healthy
echo "⏳ Waiting for dependencies to be ready..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    healthy_services=$(docker-compose ps --format json | jq -r 'select(.Health == "healthy") | .Service' 2>/dev/null | wc -l || echo "0")
    if [ "$healthy_services" -ge 3 ]; then  # postgres, redis, minio (telegramify doesn't have health check)
        echo "  ✅ Dependencies are healthy!"
        break
    else
        echo "  ⏳ Dependencies starting... (attempt $((attempt + 1))/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    fi
done

if [ $attempt -eq $max_attempts ]; then
    echo "⚠️  Warning: Some services may not be fully ready, but proceeding..."
fi

# Initialize MinIO bucket
echo "🪣 Initializing MinIO bucket..."
docker-compose up minio-init

echo "🐳 Running bot in Docker container..."
docker-compose up -d app
echo "📋 Bot is running in Docker. Use 'docker-compose logs -f app' to view logs"

echo ""
echo "🎉 Bot is running!"
echo ""
echo "🌐 Service URLs:"
echo "  MinIO Console: http://localhost:9001 (minioadmin/minioadmin)"
echo "  PostgreSQL: localhost:5432 (postgres/postgres)"
echo "  Redis: localhost:6379"
echo "  Telegramify: http://localhost:8000"
echo ""
echo "📋 Useful commands:"
echo "  View logs: docker-compose logs -f"
echo "  View app logs: docker-compose logs -f app"
echo "  Restart service: docker-compose restart [service_name]"
echo "  Stop all: docker-compose down"
echo "  Clean up: docker-compose down -v (removes volumes)"