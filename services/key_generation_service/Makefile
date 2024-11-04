# Makefile

# Load environment variables from .env
include .env

# Default target
all: start-redis run-app

# Target to start Redis
start-redis:
	@echo "Starting Redis..."
	@redis-server --port $(REDIS_PORT) &

# Target to run the Go application
run-app:
	@echo "Running the application..."
	@go run cmd/main.go

# Clean up Redis if needed
stop-redis:
	@echo "Stopping Redis..."
	@pkill redis-server

