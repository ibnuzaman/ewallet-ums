## Makefile for ewallet-ums application

# Load environment variables from .env if exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Database configuration with defaults
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= ewallet_ums
DB_SSLMODE ?= disable

# Migration settings
MIGRATION_DIR = database/migrations
DATABASE_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help build run test clean tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev/install: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest
	@echo "Installation complete!"

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/ewallet-ums main.go
	@echo "Build complete: bin/ewallet-ums"

run: ## Run the application
	@echo "Running application..."
	@go run main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Tidy complete"

dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/air-verse/air@latest"; \
		echo "Running without auto-reload..."; \
		go run main.go; \
	fi

compose-dev-up: ## Run docker-compose in development mode
	@echo "Starting docker-compose in development mode..."
	@docker-compose -f docker-compose.dev.yml up -d

compose-dev-down: ## Stop docker-compose in development mode
	@echo "Stopping docker-compose in development mode..."
	@docker-compose -f docker-compose.dev.yml down

docker-build: ## Build docker image
	@echo "Building docker image..."
	@docker build -t ewallet-ums:latest .

docker-run: ## Run docker container
	@echo "Running docker container..."
	@docker run -p 8080:8080 --env-file .env ewallet-ums:latest

lint: ## Run linter
	@echo "Running golangci-lint..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: make dev/install"; \
	fi

lint-fix: ## Run linter with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Install with: make dev/install"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi
	@if command -v gofumpt > /dev/null; then \
		gofumpt -w .; \
	fi
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

check: ## Run all checks (fmt, vet, lint, test)
	@echo "Running all checks..."
	@make fmt
	@make vet
	@make lint
	@make test
	@echo "All checks passed!"

ci: ## Run CI checks
	@echo "Running CI checks..."
	@go mod tidy
	@go fmt ./...
	@go vet ./...
	@golangci-lint run
	@go test -v -race -coverprofile=coverage.txt ./...
	@echo "CI checks passed!"

# Database commands
db-create: ## Create database
	@echo "Creating database..."
	@createdb -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME) 2>/dev/null || echo "Database might already exist"
	@echo "Database ready"

db-drop: ## Drop database
	@echo "Dropping database..."
	@dropdb -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME) 2>/dev/null || echo "Database might not exist"
	@echo "Database dropped"

db-reset: ## Reset database (drop, create, migrate)
	@echo "Resetting database..."
	@make db-drop || true
	@make db-create
	@make migrate-up
	@echo "Database reset complete"

db-status: ## Check database tables
	@echo "Checking database status..."
	@psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -c "\dt"

db-shell: ## Open database shell
	@psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME)

# Migration commands using golang-migrate
migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-up: ## Run all pending migrations
	@echo "Running migrations..."
	@echo "Database URL: $(DATABASE_URL)"
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up
	@echo "Migrations completed"

migrate-down: ## Rollback last migration
	@echo "Rolling back migration..."
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down 1
	@echo "Rollback completed"

migrate-drop: ## Drop all tables (dangerous!)
	@echo "WARNING: This will drop all tables!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" drop -f; \
		echo "All tables dropped"; \
	else \
		echo "Cancelled"; \
	fi

migrate-version: ## Show current migration version
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" version

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version parameter is required"; \
		echo "Usage: make migrate-force version=N"; \
		exit 1; \
	fi
	@echo "Forcing migration to version $(version)..."
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" force $(version)
	@echo "Migration forced to version $(version)"

migrate-goto: ## Migrate to specific version (usage: make migrate-goto version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version parameter is required"; \
		echo "Usage: make migrate-goto version=N"; \
		exit 1; \
	fi
	@echo "Migrating to version $(version)..."
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" goto $(version)
	@echo "Migrated to version $(version)"
