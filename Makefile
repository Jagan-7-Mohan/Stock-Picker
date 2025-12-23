.PHONY: build run clean test deps install

# Binary name
BINARY_NAME=stock-picker

# Build the application
build:
	go build -o $(BINARY_NAME) ./cmd/stock-picker

# Run the application
run: build
	./$(BINARY_NAME)

# Run without building
run-dev:
	go run ./cmd/stock-picker

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests (if any)
test:
	go test ./...

# Install dependencies and build
install: deps build

