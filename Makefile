.PHONY: help build run test lint clean docker-build docker-up docker-down migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@CGO_ENABLED=0 go build -ldflags='-w -s' -o bin/mpb ./cmd/main.go
	@CGO_ENABLED=0 go build -ldflags='-w -s' -o bin/migrate ./cmd/migrate.go

run: ## Run the application
	@go run ./cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@golangci-lint run

lint-install: ## Install golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

clean: ## Clean build artifacts
	@rm -rf bin/ coverage.out coverage.html
	@echo "Clean completed"

docker-build: ## Build Docker image
	@docker build -t mpb:latest .

docker-up: ## Start services with docker-compose
	@docker-compose up -d

docker-down: ## Stop services with docker-compose
	@docker-compose down

docker-logs: ## Show docker-compose logs
	@docker-compose logs -f

migrate: ## Run database migrations
	@go run ./cmd/migrate.go

migrate-up: ## Run migrations up
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations postgres "$${DSN:-postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable}" up

migrate-down: ## Run migrations down
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations postgres "$${DSN:-postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable}" down

migrate-status: ## Show migration status
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations postgres "$${DSN:-postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable}" status

deps: ## Download and verify dependencies
	@go mod download
	@go mod verify
	@go mod tidy

swagger: ## Generate swagger documentation
	@swag init -g cmd/main.go

dev: docker-up ## Start development environment
	@echo "Waiting for database..."
	@sleep 5
	@$(MAKE) migrate-up
	@echo "Development environment ready!"
	@echo "Application will be available at http://localhost:8000"
	@echo "Swagger docs at http://localhost:8000/swagger/index.html"

