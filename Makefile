.PHONY: build run test clean docker-build docker-run migrate help migrate-up migrate-down migrate-create

# Variables
BINARY_NAME=main
GO=go
GOFLAGS=-v
LINT=golangci-lint

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application in development mode"
	@echo "  test           - Run tests"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  docker-logs    - Show Docker logs"
	@echo "  migrate-up     - Run database migrations up"
	@echo "  migrate-down   - Run database migrations down"
	@echo "  migrate-create - Create a new migration"

# Build the application
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/server

# Run the application
run:
	$(GO) run ./cmd/server

# Run tests
test:
	$(GO) test -v -cover ./...

# Run tests with race detector
test-race:
	$(GO) test -v -race -cover ./...

# Run linter
lint:
	$(LINT) run ./...

# Format code
fmt:
	$(GO) fmt ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

# Docker build
docker-build:
	docker-compose -f docker/docker-compose.yml build

# Docker run
docker-run:
	docker-compose -f docker/docker-compose.yml up

# Docker stop
docker-stop:
	docker-compose -f docker/docker-compose.yml down

# Docker logs
docker-logs:
	docker-compose -f docker/docker-compose.yml logs -f

# Run database migrations up
migrate-up:
	cd cmd/migrate-cli && go run main.go up

# Run database migrations down
migrate-down:
	cd cmd/migrate-cli && go run main.go down

# Create a new migration
migrate-create:
	go run ./cmd/migrate-cli/create.go
