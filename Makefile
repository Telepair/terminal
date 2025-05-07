# Default shell
SHELL := /bin/bash

# Base directory
BASEDIR = $(shell pwd)

# Go module path (e.g., github.com/your-username/terminal)
# Attempt to auto-detect; otherwise, set it manually.
PROJECT_MODULE_PATH ?= $(shell go list -m)
ifeq ($(PROJECT_MODULE_PATH),)
    $(error Please set PROJECT_MODULE_PATH, e.g., export PROJECT_MODULE_PATH=github.com/your-username/terminal or ensure go.mod is present)
endif

# Version package path
VERSION_PKG = $(PROJECT_MODULE_PATH)/pkg/version

# Build variables
VERSION ?= $(shell git describe --tags --abbrev=0 --always --dirty 2>/dev/null || echo "v0.0.0-dev")
GIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")
GIT_COMMIT := $(shell git log -1 --pretty=format:%H 2>/dev/null || echo "unknown")
GIT_TREE_STATE := $(shell if git status --porcelain 2>/dev/null | grep -q .; then echo "dirty"; else echo "clean"; fi)
BUILD_DATE := $(shell TZ=Asia/Shanghai date +%FT%T%z)

# LDFLAGS for version injection
LDFLAGS := -X '$(VERSION_PKG).Version=$(VERSION)' \
           -X '$(VERSION_PKG).GitTag=$(GIT_TAG)' \
           -X '$(VERSION_PKG).GitCommit=$(GIT_COMMIT)' \
           -X '$(VERSION_PKG).GitBranch=$(GIT_BRANCH)' \
           -X '$(VERSION_PKG).GitTreeState=$(GIT_TREE_STATE)' \
           -X '$(VERSION_PKG).BuildDate=$(BUILD_DATE)'

# Environment: dev (default) or pro
ENV ?= dev

# Binary output directory and name
BINARY_DIR := $(BASEDIR)/bin
BINARY_NAME := terminal # Single binary name
CMD_MAIN_PATH ?= ./

# Subcommand arguments (adjust these based on your actual subcommands)
SERVER_SUBCOMMAND ?= server
AGENT_SUBCOMMAND ?= agent
TUI_SUBCOMMAND ?= tui # or cli, client, etc.

# Development mode arguments (assuming it's a flag on the server subcommand)
DEV_MODE_SERVER_ARGS ?= --dev-mode

# Build flags
BUILD_FLAGS := -v
ifeq ($(ENV),dev)
    BUILD_FLAGS += -race -gcflags="all=-N -l" # Add race detector and disable optimizations for dev
    LDFLAGS +=
else ifeq ($(ENV),pro)
    BUILD_FLAGS += -trimpath
    LDFLAGS += -s -w # Strip symbols and DWARF info for production
endif

# Tools
GOIMPORTS := $(shell go env GOPATH)/bin/goimports
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
GOLANGCI_LINT_EXISTS := $(shell command -v $(GOLANGCI_LINT) 2>/dev/null)
GOIMPORTS_EXISTS := $(shell command -v $(GOIMPORTS) 2>/dev/null)

.PHONY: all build run run-dev run-server run-agent run-tui fmt lint test test-unit cover clean help install-tools

# Default target
all: fmt lint test-unit build

# Build target
build:
	@echo "  >  Building $(BINARY_NAME) binary... (ENV=$(ENV))"
	@mkdir -p $(BINARY_DIR)
	go build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME) $(CMD_MAIN_PATH)
	@echo "  >  Binary $(BINARY_DIR)/$(BINARY_NAME) built successfully."

# Run targets
run: run-dev # Default run target, runs dev mode

run-dev: build
	@echo "  >  Running $(BINARY_NAME) in Development Mode (ENV=dev)..."
	@echo "     Executing: $(BINARY_DIR)/$(BINARY_NAME) $(SERVER_SUBCOMMAND) $(DEV_MODE_SERVER_ARGS)"
	@echo "     Access Web UI (if applicable): http://localhost:<port> (check logs for port)"
	$(BINARY_DIR)/$(BINARY_NAME) $(SERVER_SUBCOMMAND) $(DEV_MODE_SERVER_ARGS)

run-server: build
	@echo "  >  Running $(BINARY_NAME) Server..."
	@echo "     Executing: $(BINARY_DIR)/$(BINARY_NAME) $(SERVER_SUBCOMMAND)"
	$(BINARY_DIR)/$(BINARY_NAME) $(SERVER_SUBCOMMAND)

run-agent: build
	@echo "  >  Running $(BINARY_NAME) Agent..."
	@echo "     Executing: $(BINARY_DIR)/$(BINARY_NAME) $(AGENT_SUBCOMMAND)"
	# Agent might require configuration, e.g., server address
	$(BINARY_DIR)/$(BINARY_NAME) $(AGENT_SUBCOMMAND) # --server-addr=nats://localhost:4222 # Example

run-tui: build
	@echo "  >  Running $(BINARY_NAME) TUI Client..."
	@echo "     Executing: $(BINARY_DIR)/$(BINARY_NAME) $(TUI_SUBCOMMAND)"
	$(BINARY_DIR)/$(BINARY_NAME) $(TUI_SUBCOMMAND)

# Development lifecycle targets
fmt: install-tools
	@echo "  >  Formatting Go code..."
	$(GOIMPORTS) -l -w .
	@echo "  >  Running go fmt..."
	go fmt ./...

lint: install-tools
	@echo "  >  Running linters..."
	@echo "     Running go vet..."
	go vet ./...
	@echo "     Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...

test: test-unit # Alias 'test' to 'test-unit' as per plans_zh.md

test-unit:
	@echo "  >  Running unit tests..."
	go test -v -race $(shell go list ./... | grep -v /vendor/ | grep -v /test/) # Exclude vendor and e2e/integration test dirs if any

# Example target for running tests tagged as 'Example' (from original Makefile)
example:
	@echo "  >  Running example tests..."
	go test -v -run '^Example' ./...

# Benchmark target (from original Makefile)
bench:
	@echo "  >  Running benchmarks..."
	go test -v -bench=. ./...

cover:
	@echo "  >  Generating test coverage report..."
	go test ./... -v -race -short -coverprofile=.coverage.txt -covermode=atomic
	@echo "     To view HTML report: go tool cover -html=.coverage.txt"
	go tool cover -func .coverage.txt

# Clean target
clean:
	@echo "  >  Cleaning up..."
	rm -rf $(BINARY_DIR)
	rm -f .coverage.txt
	# go clean -cache -modcache # Uncomment for a more thorough clean of Go caches
	@echo "     Cleaned binary and coverage files."

# Tool installation
install-tools:
	@if [ -z "$(GOIMPORTS_EXISTS)" ]; then \
		echo "  >  Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@if [ -z "$(GOLANGCI_LINT_EXISTS)" ]; then \
		echo "  >  Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Help target
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  all                 Run fmt, lint, test-unit, and build the '$(BINARY_NAME)' binary (default)."
	@echo ""
	@echo "Build targets (ENV=dev|pro, default: dev):"
	@echo "  build               Build the single '$(BINARY_NAME)' binary."
	@echo ""
	@echo "Run targets (usually run after a build or with ENV=dev for dev builds):"
	@echo "  run-dev             Run the server in development mode (e.g., '$(BINARY_NAME) $(SERVER_SUBCOMMAND) $(DEV_MODE_SERVER_ARGS)')."
	@echo "  run-server          Run the Management Server (e.g., '$(BINARY_NAME) $(SERVER_SUBCOMMAND)')."
	@echo "  run-agent           Run the Agent (e.g., '$(BINARY_NAME) $(AGENT_SUBCOMMAND)')."
	@echo "  run-tui             Run the TUI Client (e.g., '$(BINARY_NAME) $(TUI_SUBCOMMAND)')."
	@echo ""
	@echo "Development lifecycle:"
	@echo "  fmt                 Format Go source code (goimports, go fmt)."
	@echo "  lint                Run linters (go vet, golangci-lint)."
	@echo "  test / test-unit    Run unit tests with race detector."
	@echo "  example             Run example tests."
	@echo "  bench               Run benchmarks."
	@echo "  cover               Generate and display test coverage report."
	@echo ""
	@echo "Other targets:"
	@echo "  clean               Remove built binary and coverage files."
	@echo "  install-tools       Install necessary Go development tools (goimports, golangci-lint)."
	@echo "  help                Show this help message."
	@echo ""
	@echo "Environment variables:"
	@echo "  ENV                   Set to 'pro' for production builds (e.g., make ENV=pro build). Default is 'dev'."
	@echo "  VERSION               Set a custom version string (e.g., make VERSION=v1.2.3 build)."
	@echo "  PROJECT_MODULE_PATH   Go module path, usually auto-detected (e.g., github.com/user/project)."
	@echo "  CMD_MAIN_PATH         Path to the main package for '$(BINARY_NAME)' (default: ./cmd/$(BINARY_NAME))."
	@echo "  SERVER_SUBCOMMAND     Server subcommand name (default: server)."
	@echo "  AGENT_SUBCOMMAND      Agent subcommand name (default: agent)."
	@echo "  TUI_SUBCOMMAND        TUI subcommand name (default: tui)."
	@echo "  DEV_MODE_SERVER_ARGS  Arguments for server dev mode (default: --dev-mode)."
	@echo ""