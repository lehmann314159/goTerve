.PHONY: build run test clean dev tidy

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run: build
	./bin/server

# Run in development mode with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f data/terve.db

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Create data directory
init:
	mkdir -p data
	mkdir -p bin

# Build and run
all: init build run
