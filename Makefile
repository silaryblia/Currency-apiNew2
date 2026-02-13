# Project settings
APP_NAME := currency-api
CMD_DIR  := ./cmd/api
BIN_DIR  := bin
BIN      := $(BIN_DIR)/$(APP_NAME)

GO       := go
GOFLAGS  := -mod=mod


# Tool versions
GOLANGCI_LINT_VERSION := v1.54.2
GOLANGCI_LINT_BIN := $(BIN_DIR)/golangci-lint

PROTOC_GEN_GO_VERSION := v1.33.0
PROTOC_GEN_GO_GRPC_VERSION := v1.3.0


# Default target
.DEFAULT_GOAL := help


# Helpers
.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo "  make lint            Run golangci-lint"
	@echo "  make fmt             Format code"
	@echo "  make test            Run tests"
	@echo "  make build           Build binary"
	@echo "  make run             Run application"
	@echo "  make proto           Generate protobuf files"
	@echo "  make docker          Build docker image"
	@echo "  make docker-compose  Run docker-compose"
	@echo "  make clean           Clean artifacts"
	@echo ""

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)


# Lint
.PHONY: lint
lint: $(GOLANGCI_LINT_BIN)
	@$(GOLANGCI_LINT_BIN) run ./...

$(GOLANGCI_LINT_BIN): | $(BIN_DIR)
	@echo ">> installing golangci-lint $(GOLANGCI_LINT_VERSION)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BIN_DIR) $(GOLANGCI_LINT_VERSION)


# Format
.PHONY: fmt
fmt:
	@echo ">> formatting"
	@$(GO) fmt ./...
	@command -v goimports >/dev/null 2>&1 || \
		$(GO) install golang.org/x/tools/cmd/goimports@latest
	@goimports -w .


# Test
.PHONY: test
test:
	@$(GO) test ./... -race -count=1


# Build / Run
.PHONY: build
build: | $(BIN_DIR)
	@echo ">> building $(APP_NAME)"
	@$(GO) build $(GOFLAGS) -o $(BIN) $(CMD_DIR)

.PHONY: run
run:
	@$(GO) run $(CMD_DIR)


# Protobuf
.PHONY: proto
proto:
	@command -v protoc >/dev/null 2>&1 || \
		(echo "protoc is not installed" && exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || \
		$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || \
		$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

	protoc \
		--go_out=. \
		--go-grpc_out=. \
		proto/*.proto


# Docker
.PHONY: docker
docker:
	docker build -t $(APP_NAME):latest .

.PHONY: docker-compose
docker-compose:
	docker-compose up --build


# Clean
.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)