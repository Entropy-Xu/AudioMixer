.PHONY: build run clean test install deps help

# Variables
BINARY_NAME=audio-mixer
MAIN_FILE=main.go
BUILD_DIR=build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: ./$(BINARY_NAME)"

# Build for release (with optimizations)
build-release:
	@echo "Building $(BINARY_NAME) for release..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Release build complete: ./$(BINARY_NAME)"

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Install system dependencies (macOS)
deps-macos:
	@echo "Installing PortAudio on macOS..."
	brew install portaudio
	@echo "System dependencies installed"

# Install system dependencies (Linux)
deps-linux:
	@echo "Installing PortAudio on Linux..."
	sudo apt-get update
	sudo apt-get install -y portaudio19-dev
	@echo "System dependencies installed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Lint complete"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

	@echo "Build complete for all platforms"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-release - Build with optimizations"
	@echo "  run          - Build and run the application"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install Go dependencies"
	@echo "  deps-macos   - Install system dependencies (macOS)"
	@echo "  deps-linux   - Install system dependencies (Linux)"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  build-all    - Build for all platforms"
	@echo "  help         - Show this help message"
