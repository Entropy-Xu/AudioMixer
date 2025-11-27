.PHONY: build run clean test install deps help gui build-windows build-windows-debug

# Variables
BINARY_NAME=audio-mixer
GUI_BINARY_NAME=audio-mixer-gui
MAIN_FILE=main.go
GUI_MAIN_DIR=./cmd/gui
BUILD_DIR=build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: ./$(BINARY_NAME)"

# Build GUI version
gui:
	@echo "Building $(GUI_BINARY_NAME)..."
	go build -ldflags="-s -w" -o $(GUI_BINARY_NAME) $(GUI_MAIN_DIR)
	@echo "GUI build complete: ./$(GUI_BINARY_NAME)"

# Build GUI for release (with optimizations)
build-release:
	@echo "Building $(BINARY_NAME) for release..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Building $(GUI_BINARY_NAME) for release..."
	go build -ldflags="-s -w" -o $(GUI_BINARY_NAME) $(GUI_MAIN_DIR)
	@echo "Release build complete: ./$(BINARY_NAME) and ./$(GUI_BINARY_NAME)"

# Build for Windows (with WASAPI support)
build-windows:
	@echo "Building for Windows with WASAPI support..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building CLI version..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "Building GUI version (hidden console)..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-H windowsgui -s -w" -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-windows-amd64.exe $(GUI_MAIN_DIR)
	@echo "Windows build complete:"
	@echo "  - $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe (CLI)"
	@echo "  - $(BUILD_DIR)/$(GUI_BINARY_NAME)-windows-amd64.exe (GUI with WASAPI)"

# Build for Windows (debug version with console)
build-windows-debug:
	@echo "Building Windows debug version (with console)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-windows-amd64-debug.exe $(GUI_MAIN_DIR)
	@echo "Windows debug build complete: $(BUILD_DIR)/$(GUI_BINARY_NAME)-windows-amd64-debug.exe"

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(GUI_BINARY_NAME)
	rm -f $(BINARY_NAME).exe $(GUI_BINARY_NAME).exe
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

# Build for multiple platforms (CLI only)
build-all:
	@echo "Building CLI for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

	@echo "Build complete for all platforms"

# Build GUI for multiple platforms
build-all-gui:
	@echo "Building GUI for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building GUI for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-darwin-amd64 $(GUI_MAIN_DIR)

	@echo "Building GUI for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-darwin-arm64 $(GUI_MAIN_DIR)

	@echo "Building GUI for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-linux-amd64 $(GUI_MAIN_DIR)

	@echo "Building GUI for Windows (amd64) with WASAPI..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-H windowsgui -s -w" -o $(BUILD_DIR)/$(GUI_BINARY_NAME)-windows-amd64.exe $(GUI_MAIN_DIR)

	@echo "GUI build complete for all platforms"
	@echo ""
	@echo "Windows version includes:"
	@echo "  ✓ WASAPI application audio enumeration"
	@echo "  ✓ Hidden console window"
	@echo "  ✓ go-ole COM interface support"

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Basic builds:"
	@echo "  build              - Build CLI version"
	@echo "  gui                - Build GUI version"
	@echo "  build-release      - Build both CLI and GUI with optimizations"
	@echo ""
	@echo "Windows builds:"
	@echo "  build-windows      - Build Windows version with WASAPI support"
	@echo "  build-windows-debug - Build Windows debug version (with console)"
	@echo ""
	@echo "Cross-platform builds:"
	@echo "  build-all          - Build CLI for all platforms"
	@echo "  build-all-gui      - Build GUI for all platforms (includes WASAPI on Windows)"
	@echo ""
	@echo "Development:"
	@echo "  run                - Build and run CLI"
	@echo "  clean              - Remove build artifacts"
	@echo "  test               - Run tests"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps               - Install Go dependencies"
	@echo "  deps-macos         - Install system dependencies (macOS)"
	@echo "  deps-linux         - Install system dependencies (Linux)"
	@echo ""
	@echo "Help:"
	@echo "  help               - Show this help message"
