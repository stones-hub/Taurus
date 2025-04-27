# 从 .releaserc 文件中获取版本号，如果获取不到则使用默认版本号 v1.0.0
VERSION := $(shell grep -Eo 'version:[^ ]+' .releaserc | sed -E 's/version:(.*)/\1/' || echo "v1.0.0")
# 动态获取当前目录名作为 APP_NAME 
APP_NAME := $(shell basename $(shell pwd) | tr '[:upper:]' '[:lower:]')

# 构建目录, 适用于本地化部署，不适用与docker部署
BUILD_DIR := build

# Docker parameters
# host port
HOST_PORT ?= 8080
# application in container port
CONTAINER_PORT ?= 8080
# Docker workdir
WORKDIR ?= /app
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_CONTAINER := $(APP_NAME)-container
DOCKER_NETWORK := $(APP_NAME)-network
DOCKER_VOLUME := $(APP_NAME)-config
DOCKER_LOG_VOLUME := $(APP_NAME)-logs

.PHONY: all build clean docker-build docker-run docker-stop local-run local-stop

# Default target
all: build

# Build the Go application
build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
	@echo "Build complete. Binary is located at $(BUILD_DIR)/$(APP_NAME)"

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Run the application locally
local-run: build
	@echo "Running the application locally..."
	@$(BUILD_DIR)/$(APP_NAME) -config=config/ || echo "Failed to run the application locally."

# Stop the local application (if running in the background)
local-stop:
	@echo "Stopping the local application..."
	@pkill -f "$(BUILD_DIR)/$(APP_NAME)" || echo "No local application is running."
	@echo "Local application stopped."

# Build the Docker image
docker-build:
	@echo "Removing old Docker image if it exists..."
	@docker rmi -f $(DOCKER_IMAGE) || echo "No existing image to remove."
	@echo "Building Docker image..."
	@docker build --build-arg WORKDIR=$(WORKDIR) -t $(DOCKER_IMAGE) . || echo "Failed to build Docker image."
	@echo "Docker image built: $(DOCKER_IMAGE)"

# Run the application in Docker
docker-run: docker-build
	@echo "Creating Docker network if it does not exist..."
	@docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || docker network create $(DOCKER_NETWORK)
	@echo "Checking if Docker volumes exist..."
	@if ! docker volume inspect $(DOCKER_VOLUME) >/dev/null 2>&1; then \
		echo "Creating Docker volume for config..."; \
		docker volume create $(DOCKER_VOLUME); \
		echo "Copying config files to Docker volume..."; \
		docker run --rm -v $(PWD)/config:/config -v $(DOCKER_VOLUME):$(WORKDIR)/config alpine sh -c "cp -r /config/* $(WORKDIR)/config/"; \
	fi
	@if ! docker volume inspect $(DOCKER_LOG_VOLUME) >/dev/null 2>&1; then \
		echo "Creating Docker volume for logs..."; \
		docker volume create $(DOCKER_LOG_VOLUME); \
	fi
	@echo "Running the application in Docker..."
	@docker run -d --name $(DOCKER_CONTAINER) --network $(DOCKER_NETWORK) -p ${HOST_PORT}:${CONTAINER_PORT} \
		--mount source=$(DOCKER_VOLUME),target=$(WORKDIR)/config \
		--mount source=$(DOCKER_LOG_VOLUME),target=$(WORKDIR)/logs \
		$(DOCKER_IMAGE) || echo "Failed to run the application in Docker."

# Stop the Docker container and remove the network, volumes, and image
docker-stop:
	@echo "Stopping the Docker container..."
	@docker stop $(DOCKER_CONTAINER) || echo "No running container to stop."
	@docker rm $(DOCKER_CONTAINER) || echo "No container to remove."
	@echo "Removing Docker network..."
	@docker network rm $(DOCKER_NETWORK) || echo "No network to remove."
	@echo "Removing Docker image..."
	@docker rmi $(DOCKER_IMAGE) || echo "No image to remove."
	@echo "Docker container, network, volumes, and image cleaned up."

# 使用 HOST_PORT=8090 CONTAINER_PORT=8090 WORKDIR=/myapp make docker-run 指定端口和工作目录, 默认端口为 8080， 如果不传入端口，则使用默认端口
# @echo "Removing Docker volumes..."
# @docker volume rm $(DOCKER_VOLUME) || echo "No config volume to remove."
# @docker volume rm $(DOCKER_LOG_VOLUME) || echo "No log volume to remove."