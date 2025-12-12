# Makefile for CyberArk API Command Harness

# Binary name
BINARY=cyberark

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

# Build the project
# This compiles all .go files in the current directory
build:
	@echo "Building $(BINARY)..."
	$(GOBUILD) -o $(BINARY) -v

# Install to $GOPATH/bin (makes it available in your PATH)
install:
	@echo "Installing $(BINARY)..."
	$(GOCMD) install

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)

# Run tests (when you add them)
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Format all Go files
# Go has a standard formatter - everyone uses the same style
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run the program (useful for development)
run: build
	./$(BINARY)

# Build for multiple platforms
# Go makes cross-compilation very easy
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY)-linux-amd64

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe

# Help target
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary for current platform"
	@echo "  install     - Install binary to \$$GOPATH/bin"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format all Go source files"
	@echo "  run         - Build and run the program"
	@echo "  build-all   - Build for Linux, macOS, and Windows"
	@echo "  help        - Show this help message"

.PHONY: build install clean test fmt run build-all build-linux build-darwin build-windows help

