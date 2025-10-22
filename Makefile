.PHONY: help build run test clean tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev/install:
	@echo "Installing air for development mode..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installation complete. You can now run 'make dev' to start the application in development mode."
		
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
		echo "golangci-lint not installed. Install with: brew install golangci-lint or go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

lint-fix: ## Run linter with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint or go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "Warning: goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
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
