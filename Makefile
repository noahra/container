# Makefile for container - Docker clone project

# Binary name
BINARY_NAME=container

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Source directory
SRC_DIR=./cmd/container

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)

# Build for Linux (useful if developing on different OS)
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary to $GOPATH/bin
install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/

# Run the binary (example usage)
run: build
	./$(BINARY_NAME)

# Development build with debug info
debug:
	$(GOBUILD) -gcflags="-N -l" -o $(BINARY_NAME) $(SRC_DIR)

# Build with all optimizations disabled for debugging
dev:
	$(GOBUILD) -race -o $(BINARY_NAME) $(SRC_DIR)

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the container binary"
	@echo "  build-linux - Build for Linux (cross-compile)"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  install     - Install binary to GOPATH/bin"
	@echo "  run         - Build and run the binary"
	@echo "  debug       - Build with debug symbols"
	@echo "  dev         - Build with race detection"
	@echo "  help        - Show this help message"

# Default target
.DEFAULT_GOAL := build

# Declare phony targets
.PHONY: build build-linux clean test deps install run debug dev help