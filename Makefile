.PHONY: build run clean test docker-build docker-run docker-compose-up deps restart swagger swagger-ui

# Go related variables
BINARY_NAME=url-shortener
MAIN_FILE=main.go

# Docker related variables
DOCKER_IMAGE=url-shortener
DOCKER_TAG=latest

# Default target
all: build

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Build the project
build: deps
	@echo "Building..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the project
run: build
	@echo "Running..."
	./$(BINARY_NAME)

# Run in development mode with hot reload (requires air: https://github.com/cosmtrek/air)
dev:
	@echo "Running in development mode..."
	@if ! command -v air &> /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Clean build files
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "Testing..."
	go test -v ./...

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init --parseDependency

# Run application with Swagger UI enabled
swagger-ui: swagger build
	@echo "Running with Swagger UI enabled..."
	@echo "Access Swagger UI at: http://localhost:8080/swagger/index.html"
	./$(BINARY_NAME)

# Build docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run all services with docker-compose
docker-compose-up:
	@echo "Starting all services with docker-compose..."
	docker-compose up --build -d

# Stop all services

docker-stop:
	@echo "Stopping all services..."
	docker stop $$(docker ps -q)

restart-app:
	@echo "Stopping the application..."
	docker stop $$(docker ps -q --filter "name=$(BINARY_NAME)")
	@echo "Rebuilding and starting the application..."
	docker-compose up --build -d
	@echo "Application restarted successfully!"


# Restart all services (stop and start again)
restart-all:
	@echo "Stopping all services..."
	docker stop $$(docker ps -q)
	@echo "Rebuilding and starting all services..."
	docker-compose up --build -d
	@echo "Services restarted successfully!"

# Show help
help:
	@echo "Available commands:"
	@echo "  make deps           - Download dependencies"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make swagger-ui     - Run application with Swagger UI enabled"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make docker-compose-up - Start all services with docker-compose"
	@echo "  make restart        - Stop, rebuild and restart all services with docker-compose"