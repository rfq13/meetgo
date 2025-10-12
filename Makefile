# =============================================================================
# WebRTC Meeting Application Makefile
# =============================================================================
# This Makefile provides convenient commands for development, testing, and deployment
# =============================================================================

.PHONY: help install build clean test lint format docker-up docker-down docker-logs docker-build docker-push dev prod backup restore

# Default target
.DEFAULT_GOAL := help

# Variables
COMPOSE_FILE = docker-compose.yml
COMPOSE_PROD_FILE = docker-compose.prod.yml
ENV_FILE = .env
BACKEND_DIR = backend
FRONTEND_DIR = frontend
JANUS_DIR = janus-server
NGINX_DIR = nginx

# Colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
RESET := \033[0m

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo "$(CYAN)WebRTC Meeting Application - Available Commands:$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(CYAN)Examples:$(RESET)"
	@echo "  make install          Install dependencies"
	@echo "  make dev              Start development environment"
	@echo "  make prod             Start production environment"
	@echo "  make test             Run tests"
	@echo "  make docker-build     Build Docker images"
	@echo "  make docker-up         Start services with Docker Compose"

# =============================================================================
# Setup and Installation
# =============================================================================

install: ## Install all dependencies
	@echo "$(GREEN)Installing dependencies...$(RESET)"
	@echo "$(BLUE)Installing backend dependencies...$(RESET)"
	cd $(BACKEND_DIR) && go mod download && go mod tidy
	@echo "$(BLUE)Installing frontend dependencies...$(RESET)"
	cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)Dependencies installed successfully!$(RESET)"

setup: ## Setup the project (copy .env.example to .env)
	@echo "$(GREEN)Setting up project...$(RESET)"
	@if [ ! -f $(ENV_FILE) ]; then \
		cp .env.example $(ENV_FILE); \
		echo "$(YELLOW)Created .env file from .env.example$(RESET)"; \
		echo "$(YELLOW)Please update the values in .env according to your environment$(RESET)"; \
	else \
		echo "$(YELLOW).env file already exists$(RESET)"; \
	fi
	@mkdir -p logs backups uploads
	@echo "$(GREEN)Project setup completed!$(RESET)"

# =============================================================================
# Development
# =============================================================================

dev: ## Start development environment
	@echo "$(GREEN)Starting development environment...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) --profile frontend up -d
	@echo "$(GREEN)Development environment started!$(RESET)"
	@echo "$(CYAN)Services:$(RESET)"
	@echo "  - API Server: http://localhost:8080"
	@echo "  - WebSocket Server: ws://localhost:8081"
	@echo "  - Janus Server: http://localhost:8088/janus"
	@echo "  - Frontend: http://localhost:3000"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo ""
	@echo "$(CYAN)View logs:$(RESET) make docker-logs"

dev-backend: ## Start only backend services
	@echo "$(GREEN)Starting backend services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) up -d postgres redis api websocket janus
	@echo "$(GREEN)Backend services started!$(RESET)"

dev-frontend: ## Start only frontend service
	@echo "$(GREEN)Starting frontend service...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) --profile frontend up -d frontend
	@echo "$(GREEN)Frontend service started!$(RESET)"

# =============================================================================
# Production
# =============================================================================

prod: ## Start production environment
	@echo "$(GREEN)Starting production environment...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) --profile production up -d
	@echo "$(GREEN)Production environment started!$(RESET)"
	@echo "$(CYAN)Services:$(RESET)"
	@echo "  - Nginx (HTTP): http://localhost"
	@echo "  - Nginx (HTTPS): https://localhost"
	@echo "  - API Server: http://localhost:8080"
	@echo "  - WebSocket Server: ws://localhost:8081"
	@echo "  - Janus Server: http://localhost:8088/janus"

# =============================================================================
# Building
# =============================================================================

build: ## Build all applications
	@echo "$(GREEN)Building applications...$(RESET)"
	@echo "$(BLUE)Building backend API...$(RESET)"
	cd $(BACKEND_DIR) && go build -o bin/api-server ./cmd/api
	@echo "$(BLUE)Building backend WebSocket...$(RESET)"
	cd $(BACKEND_DIR) && go build -o bin/websocket-server ./cmd/websocket
	@echo "$(BLUE)Building frontend...$(RESET)"
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)Build completed!$(RESET)"

build-backend: ## Build only backend applications
	@echo "$(GREEN)Building backend applications...$(RESET)"
	cd $(BACKEND_DIR) && go build -o bin/api-server ./cmd/api
	cd $(BACKEND_DIR) && go build -o bin/websocket-server ./cmd/websocket
	@echo "$(GREEN)Backend build completed!$(RESET)"

build-frontend: ## Build only frontend application
	@echo "$(GREEN)Building frontend application...$(RESET)"
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)Frontend build completed!$(RESET)"

# =============================================================================
# Testing
# =============================================================================

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(RESET)"
	@echo "$(BLUE)Running backend tests...$(RESET)"
	cd $(BACKEND_DIR) && go test -v -race -cover ./...
	@echo "$(BLUE)Running frontend tests...$(RESET)"
	cd $(FRONTEND_DIR) && npm test
	@echo "$(GREEN)All tests completed!$(RESET)"

test-backend: ## Run only backend tests
	@echo "$(GREEN)Running backend tests...$(RESET)"
	cd $(BACKEND_DIR) && go test -v -race -cover ./...
	@echo "$(GREEN)Backend tests completed!$(RESET)"

test-frontend: ## Run only frontend tests
	@echo "$(GREEN)Running frontend tests...$(RESET)"
	cd $(FRONTEND_DIR) && npm test
	@echo "$(GREEN)Frontend tests completed!$(RESET)"

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(RESET)"
	cd $(BACKEND_DIR) && go test -v -tags=integration ./tests/integration/...
	@echo "$(GREEN)Integration tests completed!$(RESET)"

test-e2e: ## Run end-to-end tests
	@echo "$(GREEN)Running end-to-end tests...$(RESET)"
	cd $(FRONTEND_DIR) && npm run test:e2e
	@echo "$(GREEN)E2E tests completed!$(RESET)"

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(RESET)"
	@echo "$(BLUE)Backend coverage...$(RESET)"
	cd $(BACKEND_DIR) && go test -v -race -cover ./... -coverprofile=coverage.out
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "$(BLUE)Frontend coverage...$(RESET)"
	cd $(FRONTEND_DIR) && npm run test:coverage
	@echo "$(GREEN)Coverage reports generated!$(RESET)"

# =============================================================================
# Code Quality
# =============================================================================

lint: ## Run linting
	@echo "$(GREEN)Running linters...$(RESET)"
	@echo "$(BLUE)Linting backend...$(RESET)"
	cd $(BACKEND_DIR) && golangci-lint run
	@echo "$(BLUE)Linting frontend...$(RESET)"
	cd $(FRONTEND_DIR) && npm run lint
	@echo "$(GREEN)Linting completed!$(RESET)"

lint-backend: ## Run backend linting
	@echo "$(GREEN)Linting backend...$(RESET)"
	cd $(BACKEND_DIR) && golangci-lint run
	@echo "$(GREEN)Backend linting completed!$(RESET)"

lint-frontend: ## Run frontend linting
	@echo "$(GREEN)Linting frontend...$(RESET)"
	cd $(FRONTEND_DIR) && npm run lint
	@echo "$(GREEN)Frontend linting completed!$(RESET)"

format: ## Format code
	@echo "$(GREEN)Formatting code...$(RESET)"
	@echo "$(BLUE)Formatting backend...$(RESET)"
	cd $(BACKEND_DIR) && go fmt ./...
	@echo "$(BLUE)Formatting frontend...$(RESET)"
	cd $(FRONTEND_DIR) && npm run format
	@echo "$(GREEN)Code formatting completed!$(RESET)"

format-backend: ## Format backend code
	@echo "$(GREEN)Formatting backend code...$(RESET)"
	cd $(BACKEND_DIR) && go fmt ./...
	@echo "$(GREEN)Backend code formatted!$(RESET)"

format-frontend: ## Format frontend code
	@echo "$(GREEN)Formatting frontend code...$(RESET)"
	cd $(FRONTEND_DIR) && npm run format
	@echo "$(GREEN)Frontend code formatted!$(RESET)"

# =============================================================================
# Docker Commands
# =============================================================================

docker-build: ## Build Docker images
	@echo "$(GREEN)Building Docker images...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) build
	@echo "$(GREEN)Docker images built!$(RESET)"

docker-build-backend: ## Build backend Docker images
	@echo "$(GREEN)Building backend Docker images...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) build api websocket janus
	@echo "$(GREEN)Backend Docker images built!$(RESET)"

docker-build-frontend: ## Build frontend Docker image
	@echo "$(GREEN)Building frontend Docker image...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) build frontend
	@echo "$(GREEN)Frontend Docker image built!$(RESET)"

docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)Starting services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started!$(RESET)"

docker-down: ## Stop all services
	@echo "$(GREEN)Stopping services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) down
	@echo "$(GREEN)Services stopped!$(RESET)"

docker-restart: ## Restart all services
	@echo "$(GREEN)Restarting services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) restart
	@echo "$(GREEN)Services restarted!$(RESET)"

docker-logs: ## Show Docker logs
	@echo "$(GREEN)Showing logs...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) logs -f

docker-logs-api: ## Show API logs
	@echo "$(GREEN)Showing API logs...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) logs -f api

docker-logs-websocket: ## Show WebSocket logs
	@echo "$(GREEN)Showing WebSocket logs...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) logs -f websocket

docker-logs-janus: ## Show Janus logs
	@echo "$(GREEN)Showing Janus logs...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) logs -f janus

docker-clean: ## Clean Docker resources
	@echo "$(GREEN)Cleaning Docker resources...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker resources cleaned!$(RESET)"

docker-push: ## Push Docker images to registry
	@echo "$(GREEN)Pushing Docker images...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) push
	@echo "$(GREEN)Docker images pushed!$(RESET)"

# =============================================================================
# Database Commands
# =============================================================================

db-migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) exec api go run ./cmd/migrate/main.go up
	@echo "$(GREEN)Migrations completed!$(RESET)"

db-migrate-down: ## Rollback database migrations
	@echo "$(GREEN)Rolling back database migrations...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) exec api go run ./cmd/migrate/main.go down
	@echo "$(GREEN)Rollback completed!$(RESET)"

db-reset: ## Reset database
	@echo "$(GREEN)Resetting database...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) down postgres
	docker volume rm webrtc-meeting_postgres_data || true
	docker-compose -f $(COMPOSE_FILE) up -d postgres
	@echo "$(GREEN)Database reset completed!$(RESET)"

db-seed: ## Seed database with sample data
	@echo "$(GREEN)Seeding database...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) exec api go run ./cmd/seed/main.go
	@echo "$(GREEN)Database seeded!$(RESET)"

db-backup: ## Backup database
	@echo "$(GREEN)Backing up database...$(RESET)"
	mkdir -p backups
	docker-compose -f $(COMPOSE_FILE) exec postgres pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Database backup completed!$(RESET)"

db-restore: ## Restore database from backup
	@echo "$(YELLOW)Usage: make db-restore BACKUP_FILE=backup_file.sql$(RESET)"

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Clean build artifacts and temporary files
	@echo "$(GREEN)Cleaning up...$(RESET)"
	@echo "$(BLUE)Cleaning backend...$(RESET)"
	cd $(BACKEND_DIR) && go clean && rm -rf bin/
	@echo "$(BLUE)Cleaning frontend...$(RESET)"
	cd $(FRONTEND_DIR) && rm -rf dist/ node_modules/.cache/
	@echo "$(BLUE)Cleaning logs...$(RESET)"
	rm -rf logs/*
	@echo "$(GREEN)Cleanup completed!$(RESET)"

clean-all: ## Clean everything including Docker resources
	@echo "$(GREEN)Deep cleaning...$(RESET)"
	$(MAKE) clean
	$(MAKE) docker-clean
	@echo "$(GREEN)Deep cleanup completed!$(RESET)"

# =============================================================================
# Monitoring and Health Checks
# =============================================================================

health: ## Check health of all services
	@echo "$(GREEN)Checking service health...$(RESET)"
	@echo "$(BLUE)API Server:$(RESET)"
	@curl -f http://localhost:8080/health || echo "$(RED)API Server is down$(RESET)"
	@echo "$(BLUE)WebSocket Server:$(RESET)"
	@curl -f http://localhost:8081/health || echo "$(RED)WebSocket Server is down$(RESET)"
	@echo "$(BLUE)Janus Server:$(RESET)"
	@curl -f http://localhost:8088/janus/info || echo "$(RED)Janus Server is down$(RESET)"
	@echo "$(GREEN)Health check completed!$(RESET)"

status: ## Show status of all services
	@echo "$(GREEN)Service status:$(RESET)"
	docker-compose -f $(COMPOSE_FILE) ps

# =============================================================================
# Utilities
# =============================================================================

shell-api: ## Open shell in API container
	docker-compose -f $(COMPOSE_FILE) exec api sh

shell-websocket: ## Open shell in WebSocket container
	docker-compose -f $(COMPOSE_FILE) exec websocket sh

shell-janus: ## Open shell in Janus container
	docker-compose -f $(COMPOSE_FILE) exec janus sh

shell-postgres: ## Open shell in PostgreSQL container
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

shell-redis: ## Open shell in Redis container
	docker-compose -f $(COMPOSE_FILE) exec redis redis-cli

logs-api: ## Follow API logs
	docker-compose -f $(COMPOSE_FILE) logs -f api

logs-websocket: ## Follow WebSocket logs
	docker-compose -f $(COMPOSE_FILE) logs -f websocket

logs-janus: ## Follow Janus logs
	docker-compose -f $(COMPOSE_FILE) logs -f janus

logs-postgres: ## Follow PostgreSQL logs
	docker-compose -f $(COMPOSE_FILE) logs -f postgres

logs-redis: ## Follow Redis logs
	docker-compose -f $(COMPOSE_FILE) logs -f redis

# =============================================================================
# Deployment
# =============================================================================

deploy-staging: ## Deploy to staging environment
	@echo "$(GREEN)Deploying to staging...$(RESET)"
	# Add staging deployment commands here
	@echo "$(GREEN)Staging deployment completed!$(RESET)"

deploy-production: ## Deploy to production environment
	@echo "$(GREEN)Deploying to production...$(RESET)"
	# Add production deployment commands here
	@echo "$(GREEN)Production deployment completed!$(RESET)"

# =============================================================================
# Version and Info
# =============================================================================

version: ## Show version information
	@echo "$(CYAN)WebRTC Meeting Application$(RESET)"
	@echo "$(BLUE)Version: 1.0.0$(RESET)"
	@echo "$(BLUE)Git: $(shell git rev-parse --short HEAD)$(RESET)"
	@echo "$(BLUE)Branch: $(shell git rev-parse --abbrev-ref HEAD)$(RESET)"

info: ## Show project information
	@echo "$(CYAN)WebRTC Meeting Application Information$(RESET)"
	@echo ""
	@echo "$(BLUE)Project Structure:$(RESET)"
	@echo "  - Backend: $(BACKEND_DIR)"
	@echo "  - Frontend: $(FRONTEND_DIR)"
	@echo "  - Janus Server: $(JANUS_DIR)"
	@echo "  - Docker Compose: $(COMPOSE_FILE)"
	@echo ""
	@echo "$(BLUE)Development URLs:$(RESET)"
	@echo "  - API Server: http://localhost:8080"
	@echo "  - WebSocket Server: ws://localhost:8081"
	@echo "  - Janus Server: http://localhost:8088/janus"
	@echo "  - Frontend: http://localhost:3000"
	@echo ""
	@echo "$(BLUE)Database:$(RESET)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"