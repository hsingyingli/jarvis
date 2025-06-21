.PHONY: build test lint run clean install-deps

# Build the application
build:
	go build -o bin/jarvis ./cmd/jarvis

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Run the application
run:
	go run ./cmd/jarvis

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
install-deps:
	go mod download
	go mod tidy

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL not set. Using default: postgresql://postgres:password@localhost:5432/jarvis?sslmode=disable"; \
		export DATABASE_URL="postgresql://postgres:password@localhost:5432/jarvis?sslmode=disable"; \
	fi; \
	psql $(DATABASE_URL) -f migrations/001_initial_schema.sql

# Development setup
dev-setup: install-deps migrate-up
	cp .env.example .env
	@echo "Development setup complete. Please edit .env file with your configuration."

# Docker build
docker-build:
	docker build -t jarvis:latest .

# Docker run
docker-run:
	docker run -p 8080:8080 jarvis:latest

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Security scan
security:
	gosec ./...

# Full check (format, vet, lint, test)
check: fmt vet lint test

# Quick build and run
dev: build run