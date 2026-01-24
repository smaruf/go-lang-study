# Makefile for 3-weeks-plan Task Management API

.PHONY: help run build test test-verbose clean install setup

# Default target
help:
	@echo "Available targets:"
	@echo "  setup         - Install dependencies and setup environment"
	@echo "  run           - Run the application"
	@echo "  build         - Build the application binary"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts and database"
	@echo "  install       - Install dependencies"

# Setup environment
setup: install
	@echo "Setting up environment..."
	@cp -n .env.example .env || true
	@echo "Environment setup complete!"

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run the application
run:
	@echo "Starting application..."
	@go run cmd/api/main.go

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/api cmd/api/main.go
	@echo "Binary built: bin/api"

# Run tests
test:
	@echo "Running tests..."
	@go test ./tests/...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./tests/...
	@go test -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts and database
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf data/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Run the built binary
run-bin: build
	@echo "Running binary..."
	@./bin/api
