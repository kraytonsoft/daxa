# Binaries
BIN_DIR := ./bin
CLI_BIN := $(BIN_DIR)/daxa
RUNTIME_BIN := $(BIN_DIR)/daxa-runtime

# Entrypoints
CLI_MAIN := ./cli
RUNTIME_MAIN := ./runtime

# Docker
DOCKER_IMAGE := daxagrid/runtime
DOCKER_TAG := latest

.PHONY: all cli runtime docker-build docker-run clean run-cli run-runtime

all: cli runtime

cli:
	@echo "🔧 Building CLI..."
	@mkdir -p $(BIN_DIR)
	go build -o $(CLI_BIN) $(CLI_MAIN)

runtime:
	@echo "🔧 Building Runtime..."
	@mkdir -p $(BIN_DIR)
	go build -o $(RUNTIME_BIN) $(RUNTIME_MAIN)

run-cli: cli
	@echo "🚀 Running CLI..."
	@$(CLI_BIN)

run-runtime: runtime
	@echo "🚀 Running Runtime on localhost:36365"
	@$(RUNTIME_BIN)

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -f ./docker/Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "🐳 Running Daxa runtime in Docker..."
	docker run --rm -p 443:443 -p 8080:8080 -p 36365:36365 $(DOCKER_IMAGE):$(DOCKER_TAG)

clean:
	@echo "🧹 Cleaning binaries..."
	@rm -rf $(BIN_DIR)
