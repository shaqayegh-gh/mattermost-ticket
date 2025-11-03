# Build information
GO ?= $(shell which go)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Plugin information
PLUGIN_NAME := com.github.mattermost-ticket-plugin
PLUGIN_VERSION := 1.0.0

# Build directories
DIST_DIR := dist
BUILD_DIR := build

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
BUILD_FLAGS := -ldflags "-X main.Version=$(PLUGIN_VERSION)"

# Output binary name
BINARY_NAME := plugin-$(GOOS)-$(GOARCH)
ifeq ($(GOOS),windows)
	BINARY_NAME := $(BINARY_NAME).exe
endif

.PHONY: all build clean test deps check-style dist

all: check-style test build

# Build the plugin
build:
	@echo Building plugin for $(GOOS)-$(GOARCH)
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME) ./server

# Clean build artifacts
clean:
	@echo Cleaning build artifacts
	$(GOCLEAN)
	rm -rf $(DIST_DIR)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo Running tests
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo Downloading dependencies
	$(GOMOD) download
	$(GOMOD) tidy

# Check code style
check-style:
	@echo Checking code style
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping style check"; \
	fi

# Build for all platforms
dist: clean
	@echo Building for all platforms
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(MAKE) build
	GOOS=darwin GOARCH=amd64 $(MAKE) build
	GOOS=windows GOARCH=amd64 $(MAKE) build

# Install the plugin
install: build
	@echo Installing plugin
	@if [ -f $(DIST_DIR)/$(BINARY_NAME) ]; then \
		echo "Plugin built successfully: $(DIST_DIR)/$(BINARY_NAME)"; \
	else \
		echo "Plugin build failed"; \
		exit 1; \
	fi

# Development build with live reload
dev:
	@echo Starting development build
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found, install with: go install github.com/cosmtrek/air@latest"; \
		$(MAKE) build; \
	fi

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the plugin for current platform"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Download dependencies"
	@echo "  check-style - Check code style with golangci-lint"
	@echo "  dist       - Build for all platforms"
	@echo "  install    - Build and show installation info"
	@echo "  dev        - Development build with live reload"
	@echo "  help       - Show this help message"
