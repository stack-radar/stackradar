# StackRadar
# Cross-platform build automation

.PHONY: all build build-linux build-mac build-windows clean install test help

# Version and build info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Binary names
BINARY_NAME := stackradar
BUILD_DIR := build

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Default target
all: clean test build

## help: Display this help message
help:
	@echo "TechStack Detector - Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^##/ {sub(/^## /, "", $$0); printf "  %-20s %s\n", $$2, substr($$0, index($$0, $$3))}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make build           # Build for current platform"
	@echo "  make build-all       # Build for all platforms"
	@echo "  make install         # Install to GOPATH/bin"
	@echo "  make test            # Run tests"

## build: Build for current platform
build:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)"

## build-linux: Build for Linux (amd64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

## build-linux-arm64: Build for Linux (arm64)
build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64"

## build-mac: Build for macOS (amd64 and arm64)
build-mac: build-mac-amd64 build-mac-arm64

build-mac-amd64:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64"

build-mac-arm64:
	@echo "Building for macOS (arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64"

## build-windows: Build for Windows (amd64)
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

## build-all: Build for all supported platforms
build-all: build-linux build-linux-arm64 build-mac build-windows
	@echo ""
	@echo "✓ All binaries built successfully!"
	@echo ""
	@ls -lh $(BUILD_DIR)

## install: Install binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(LDFLAGS) .
	@echo "✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests completed"

## test-coverage: Run tests with coverage report
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)
	@echo "✓ Cleaned"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify
	@echo "✓ Dependencies ready"

## tidy: Tidy go.mod and go.sum
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "✓ Dependencies tidied"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✓ Code formatted"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "✓ go vet passed"

## lint: Run linters (requires golangci-lint)
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install: https://golangci-lint.run/usage/install/"; exit 1)
	golangci-lint run ./...
	@echo "✓ Linting passed"

## check: Run fmt, vet, and test
check: fmt vet test

## security: Run security scans (trivy, govulncheck, gosec)
security:
	@echo "Running security scans..."
	@which trivy > /dev/null || (echo "⚠️  trivy not found. Install: brew install trivy (macOS) or see README"; exit 1)
	@echo "Running Trivy..."
	trivy fs .
	@echo ""
	@which govulncheck > /dev/null || (echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@latest)
	@echo "Running govulncheck..."
	govulncheck ./...
	@echo ""
	@which gosec > /dev/null || (echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest)
	@echo "Running GoSec..."
	gosec ./...
	@echo "✓ Security scans completed"

## release: Create release binaries with checksums
release: clean build-all
	@echo "Creating release artifacts..."
	@cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "✓ Release artifacts ready in $(BUILD_DIR)/"
	@cat $(BUILD_DIR)/checksums.txt

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t stackradar:$(VERSION) .
	docker tag stackradar:$(VERSION) stackradar:latest
	@echo "✓ Docker image built: stackradar:$(VERSION)"

## run: Build and run the application
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME) --help

## version: Display version information
version:
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
