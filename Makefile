# Variables
PROJECT_NAME=ccrun
BINARY_DIR=bin
CMD_DIR=cmd/ccrun
MAIN_FILE=$(CMD_DIR)/main.go
BINARY_PATH=$(BINARY_DIR)/$(PROJECT_NAME)

# Default target
all: build

# Create bin directory if it doesn't exist
$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

# Build the binary
build: $(BINARY_DIR)
	go build -o $(BINARY_PATH) $(MAIN_FILE)

# Run the program with "echo hello" arguments
run: build
	$(BINARY_PATH) echo hello

# Clean up the binary
clean:
	rm -f $(BINARY_PATH)

# Test target: build, run with echo hello, then clean
test: build run clean

# Declare phony targets (targets that don't represent files)
.PHONY: all build run clean test