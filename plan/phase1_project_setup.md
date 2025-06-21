# Phase 1 - Project Setup & Infrastructure

## 🎯 Goal
Establish a robust, scalable foundation with proper development workflows and monitoring capabilities.

## ✅ Success Criteria
- Project builds without errors using `go build`
- All tests pass with `go test ./...`
- CI/CD pipeline runs successfully on commit
- Logs are properly formatted and stored
- Database migrations run successfully
- Configuration loads from environment variables

## 📋 Execution Plan

### Step 1: Initialize Go Project Structure
**Duration**: 4-6 hours

#### 1.1 Create Go Module and Basic Structure
```bash
# Initialize Go module
go mod init github.com/hsingyingli/jarvis

# Create basic directory structure
mkdir -p {cmd,internal,pkg,api,web,docs,scripts,configs,migrations}
mkdir -p {internal/{core,services,handlers,middleware,models,repository},pkg/{logger,config,database,utils}}
```

#### 1.2 Create Main Application Entry Point
**File**: `cmd/jarvis/main.go`
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/hsingyingli/jarvis/internal/core"
    "github.com/hsingyingli/jarvis/pkg/config"
    "github.com/hsingyingli/jarvis/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }
    
    // Initialize logger
    logger := logger.New(cfg.LogLevel)
    
    // Initialize core application
    app, err := core.NewApp(cfg, logger)
    if err != nil {
        logger.Fatal("Failed to initialize application:", err)
    }
    
    // Start application
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        logger.Info("Shutting down...")
        cancel()
    }()
    
    if err := app.Run(ctx); err != nil {
        logger.Fatal("Application error:", err)
    }
}
```

#### 1.3 Create Core Application Structure
**File**: `internal/core/app.go`
```go
package core

import (
    "context"
    "fmt"
    
    "github.com/hsingyingli/jarvis/pkg/config"
    "github.com/hsingyingli/jarvis/pkg/logger"
)

type App struct {
    config *config.Config
    logger logger.Logger
}

func NewApp(cfg *config.Config, log logger.Logger) (*App, error) {
    return &App{
        config: cfg,
        logger: log,
    }, nil
}

func (a *App) Run(ctx context.Context) error {
    a.logger.Info("Starting Jarvis AI Assistant...")
    
    // TODO: Initialize services
    // TODO: Start HTTP server
    // TODO: Start WebSocket server
    
    <-ctx.Done()
    a.logger.Info("Application stopped")
    return nil
}
```

### Step 2: Set Up Configuration Management
**Duration**: 2-3 hours

#### 2.1 Create Configuration Package
**File**: `pkg/config/config.go`
```go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    // Server configuration
    HTTPPort     string `json:"http_port"`
    HTTPHost     string `json:"http_host"`
    
    // Database configuration
    DatabaseURL  string `json:"database_url"`
    DatabaseType string `json:"database_type"` // sqlite, postgres
    
    // AI configuration
    OpenAIAPIKey string `json:"openai_api_key"`
    ModelName    string `json:"model_name"`
    
    // Logging configuration
    LogLevel     string `json:"log_level"`
    
    // Security
    JWTSecret    string `json:"jwt_secret"`
}

func Load() (*Config, error) {
    cfg := &Config{
        HTTPPort:     getEnv("HTTP_PORT", "8080"),
        HTTPHost:     getEnv("HTTP_HOST", "localhost"),
        DatabaseURL:  getEnv("DATABASE_URL", "jarvis.db"),
        DatabaseType: getEnv("DATABASE_TYPE", "sqlite"),
        OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
        ModelName:    getEnv("MODEL_NAME", "gpt-3.5-turbo"),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
        JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

#### 2.2 Create Environment Configuration
**File**: `.env.example`
```bash
# Server Configuration
HTTP_PORT=8080
HTTP_HOST=localhost

# Database Configuration
DATABASE_URL=jarvis.db
DATABASE_TYPE=sqlite

# AI Configuration
OPENAI_API_KEY=your_openai_api_key_here
MODEL_NAME=gpt-3.5-turbo

# Logging Configuration
LOG_LEVEL=info

# Security
JWT_SECRET=your_jwt_secret_here
```

### Step 3: Set Up Logging Framework
**Duration**: 2-3 hours

#### 3.1 Create Logger Package
**File**: `pkg/logger/logger.go`
```go
package logger

import (
    "log/slog"
    "os"
)

type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Fatal(msg string, args ...interface{})
}

type slogLogger struct {
    logger *slog.Logger
}

func New(level string) Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    
    return &slogLogger{logger: logger}
}

func (l *slogLogger) Debug(msg string, args ...interface{}) {
    l.logger.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...interface{}) {
    l.logger.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...interface{}) {
    l.logger.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...interface{}) {
    l.logger.Error(msg, args...)
}

func (l *slogLogger) Fatal(msg string, args ...interface{}) {
    l.logger.Error(msg, args...)
    os.Exit(1)
}
```

### Step 4: Set Up Database Schema (PostgreSQL + Redis)
**Duration**: 3-4 hours

#### 4.1 Create Database Package
**File**: `pkg/database/database.go`
```go
package database

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
    "github.com/hsingyingli/jarvis/pkg/config"
)

type DB struct {
    Postgres *pgxpool.Pool
    Redis    *redis.Client
}

func New(cfg *config.Config) (*DB, error) {
    // Connect to PostgreSQL
    var pgConnString string
    if cfg.DatabaseURL != "" {
        pgConnString = cfg.DatabaseURL
    } else {
        pgConnString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
            cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
    }

    pgPool, err := pgxpool.New(context.Background(), pgConnString)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }

    // Connect to Redis
    var redisOpt *redis.Options
    if cfg.RedisURL != "" {
        redisOpt, err = redis.ParseURL(cfg.RedisURL)
        if err != nil {
            pgPool.Close()
            return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
        }
    } else {
        redisOpt = &redis.Options{
            Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
            Password: cfg.RedisPassword,
            DB:       cfg.RedisDB,
        }
    }

    redisClient := redis.NewClient(redisOpt)

    // Test connections
    if err := pgPool.Ping(context.Background()); err != nil {
        pgPool.Close()
        return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
    }

    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        pgPool.Close()
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &DB{
        Postgres: pgPool,
        Redis:    redisClient,
    }, nil
}
```

#### 4.2 Create Initial Migration (PostgreSQL)
**File**: `migrations/001_initial_schema.sql`
```sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- AI memory/context table for semantic search
CREATE TABLE IF NOT EXISTS ai_memories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    context_type VARCHAR(100) NOT NULL,
    importance_score FLOAT DEFAULT 0.5,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes and triggers for performance
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_ai_memories_user_id ON ai_memories(user_id);
```

### Step 5: Implement Basic Error Handling
**Duration**: 2-3 hours

#### 5.1 Create Error Package
**File**: `pkg/utils/errors.go`
```go
package utils

import (
    "fmt"
    "net/http"
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
    }
    return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// Common error constructors
func BadRequestError(message string, err error) *AppError {
    return NewAppError(http.StatusBadRequest, message, err)
}

func UnauthorizedError(message string, err error) *AppError {
    return NewAppError(http.StatusUnauthorized, message, err)
}

func NotFoundError(message string, err error) *AppError {
    return NewAppError(http.StatusNotFound, message, err)
}

func InternalServerError(message string, err error) *AppError {
    return NewAppError(http.StatusInternalServerError, message, err)
}
```

### Step 6: Set Up CI/CD Pipeline
**Duration**: 3-4 hours

#### 6.1 Create GitHub Actions Workflow
**File**: `.github/workflows/ci.yml`
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Build
      run: go build -v ./cmd/jarvis
    
    - name: Run security scan
      uses: securecodewarrior/github-action-add-sarif@v1
      with:
        sarif-file: 'gosec.sarif'
      continue-on-error: true

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Build binary
      run: |
        CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o jarvis ./cmd/jarvis
    
    - name: Upload binary
      uses: actions/upload-artifact@v3
      with:
        name: jarvis-binary
        path: jarvis
```

### Step 7: Create Basic Tests
**Duration**: 2-3 hours

#### 7.1 Create Test Files
**File**: `pkg/config/config_test.go`
```go
package config

import (
    "os"
    "testing"
)

func TestLoad(t *testing.T) {
    // Set test environment variables
    os.Setenv("HTTP_PORT", "9090")
    os.Setenv("LOG_LEVEL", "debug")
    
    defer func() {
        os.Unsetenv("HTTP_PORT")
        os.Unsetenv("LOG_LEVEL")
    }()
    
    cfg, err := Load()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if cfg.HTTPPort != "9090" {
        t.Errorf("Expected HTTP_PORT to be 9090, got %s", cfg.HTTPPort)
    }
    
    if cfg.LogLevel != "debug" {
        t.Errorf("Expected LOG_LEVEL to be debug, got %s", cfg.LogLevel)
    }
}
```

**File**: `pkg/database/database_test.go`
```go
package database

import (
    "testing"
)

func TestNew(t *testing.T) {
    db, err := New("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    defer db.Close()
    
    if db.DB == nil {
        t.Error("Expected database connection, got nil")
    }
}

func TestNewUnsupportedType(t *testing.T) {
    _, err := New("unsupported", "test.db")
    if err == nil {
        t.Error("Expected error for unsupported database type")
    }
}
```

### Step 8: Create Development Scripts
**Duration**: 1-2 hours

#### 8.1 Create Makefile
**File**: `Makefile`
```makefile
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
```

## 🚀 Execution Checklist

### Week 1
- [ ] Initialize Go project structure and modules
- [ ] Create main application entry point
- [ ] Set up configuration management system
- [ ] Implement logging framework
- [ ] Create basic error handling utilities

### Week 2
- [ ] Set up database schema and migrations
- [ ] Create CI/CD pipeline with GitHub Actions
- [ ] Write initial unit tests
- [ ] Create development scripts and Makefile
- [ ] Set up Docker configuration (optional)

## 📦 Dependencies to Add

```bash
go get github.com/jackc/pgx/v5
go get github.com/redis/go-redis/v9
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get github.com/golang-jwt/jwt/v5
```

## ⚡ Quick Start Commands

```bash
# Initialize the project
make dev-setup

# Run tests
make test

# Build the application
make build

# Run the application
make run
```

## 🔍 Verification Steps

1. **Build Test**: `go build ./cmd/jarvis` should complete without errors
2. **Test Suite**: `go test ./...` should pass all tests
3. **Configuration**: Environment variables should load correctly
4. **PostgreSQL**: Database migrations should run successfully
5. **Redis**: Cache connection should work
6. **Logging**: Structured logs should output to stdout
7. **CI/CD**: GitHub Actions should run on push/PR

## 📋 Deliverables

- ✅ Complete Go project structure
- ✅ Working configuration management
- ✅ Structured logging system
- ✅ PostgreSQL database schema and migrations
- ✅ Redis cache integration
- ✅ Error handling framework
- ✅ CI/CD pipeline
- ✅ Test suite foundation
- ✅ Development tooling and scripts