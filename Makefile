# Makefile for cleanup utility

BINARY_NAME=cleanup
SOURCE_FILE=cleanup.go
INSTALL_PATH=$(HOME)/bin

.PHONY: all build install clean help

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(SOURCE_FILE)
	@echo "Build complete: $(BINARY_NAME)"

# Install the binary to $HOME/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/
	chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete. You can now use '$(BINARY_NAME)' from anywhere."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete."

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstall complete."

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  install   - Build and install to $(INSTALL_PATH)"
	@echo "  clean     - Remove build artifacts"
	@echo "  uninstall - Remove installed binary"
	@echo "  help      - Show this help message"