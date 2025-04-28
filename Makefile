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

# Docker image name
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
# Docker container name
DOCKER_CONTAINER := $(APP_NAME)-container
# Docker network name
DOCKER_NETWORK := $(APP_NAME)-network
# Docker volume name
DOCKER_CONFIG_VOLUME := $(APP_NAME)_config_data
# Docker log volume name
DOCKER_LOG_VOLUME := $(APP_NAME)_log_data

# 定义颜色
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
CYAN := \033[36m
RESET := \033[0m

# 定义分割符
SEPARATOR := $(CYAN)========================================$(RESET)

.PHONY: all build clean docker-build docker-run docker-stop local-run local-stop

# Default target
all: build

# Build the Go application
build:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Building the application...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
	@echo -e "$(GREEN)Build complete. Binary is located at $(BUILD_DIR)/$(APP_NAME)$(RESET)"
	@echo -e "$(SEPARATOR)"

# Clean up build artifacts
clean:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(YELLOW)Cleaning up...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@echo -e "$(GREEN)Clean complete.$(RESET)"
	@echo -e "$(SEPARATOR)"

# Run the application locally
local-run: build
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Running the application locally...$(RESET)"
	@$(BUILD_DIR)/$(APP_NAME) -config=config/ || echo -e "$(RED)Failed to run the application locally.$(RESET)"
	@echo -e "$(SEPARATOR)"

# Stop the local application (if running in the background)
local-stop:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(YELLOW)Stopping the local application...$(RESET)"
	@pkill -f "$(BUILD_DIR)/$(APP_NAME)" || echo -e "$(RED)No local application is running.$(RESET)"
	@echo -e "$(GREEN)Local application stopped.$(RESET)"
	@echo -e "$(SEPARATOR)"

# Build the Docker image
docker-build:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Removing old Docker image if it exists...$(RESET)"
	@docker rmi -f $(DOCKER_IMAGE) || echo -e "$(YELLOW)No existing image to remove.$(RESET)"
	@echo -e "$(BLUE)Building Docker image...$(RESET)"
	@docker build --build-arg WORKDIR=$(WORKDIR) -t $(DOCKER_IMAGE) . || echo -e "$(RED)Failed to build Docker image.$(RESET)"
	@echo -e "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(RESET)"
	@echo -e "$(SEPARATOR)"

# Run the application in Docker
docker-run: docker-build
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Creating Docker network if it does not exist...$(RESET)"
	@docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || docker network create $(DOCKER_NETWORK)
	@echo -e "$(BLUE)Checking if Docker volumes exist...$(RESET)"
	@if ! docker volume inspect $(DOCKER_CONFIG_VOLUME) >/dev/null 2>&1; then \
		echo -e "$(YELLOW)Creating Docker volume for config...$(RESET)"; \
		docker volume create $(DOCKER_CONFIG_VOLUME); \
		echo -e "$(BLUE)Copying config files to Docker volume...$(RESET)"; \
		docker run --rm -v $(PWD)/config:/config -v $(DOCKER_CONFIG_VOLUME):$(WORKDIR)/config alpine sh -c "cp -r /config/* $(WORKDIR)/config/"; \
	fi
	@if ! docker volume inspect $(DOCKER_LOG_VOLUME) >/dev/null 2>&1; then \
		echo -e "$(YELLOW)Creating Docker volume for logs...$(RESET)"; \
		docker volume create $(DOCKER_LOG_VOLUME); \
	fi
	@echo -e "$(BLUE)Running the application in Docker...$(RESET)"
	@docker run -d --name $(DOCKER_CONTAINER) --network $(DOCKER_NETWORK) -p ${HOST_PORT}:${CONTAINER_PORT} \
		--mount source=$(DOCKER_CONFIG_VOLUME),target=$(WORKDIR)/config \
		--mount source=$(DOCKER_LOG_VOLUME),target=$(WORKDIR)/logs \
		$(DOCKER_IMAGE) || echo -e "$(RED)Failed to run the application in Docker.$(RESET)"
	@echo -e "$(SEPARATOR)"

# Stop the Docker container and remove the network and image
docker-stop:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(YELLOW)Stopping the Docker container...$(RESET)"
	@docker stop $(DOCKER_CONTAINER) || echo -e "$(RED)No running container to stop.$(RESET)"
	@docker rm $(DOCKER_CONTAINER) || echo -e "$(RED)No container to remove.$(RESET)"
	@echo -e "$(YELLOW)Removing Docker network...$(RESET)"
	@docker network rm $(DOCKER_NETWORK) || echo -e "$(RED)No network to remove.$(RESET)"
	@echo -e "$(YELLOW)Removing Docker image...$(RESET)"
	@docker rmi $(DOCKER_IMAGE) || echo -e "$(RED)No image to remove.$(RESET)"
	@echo -e "$(GREEN)Docker container, network and image cleaned up.$(RESET)"
	@echo -e "$(SEPARATOR)"

# 使用 HOST_PORT=8090 CONTAINER_PORT=8090 WORKDIR=/myapp make docker-run 指定端口和工作目录, 默认端口为 8080， 如果不传入端口，则使用默认端口
# @echo "Removing Docker volumes..."
# @docker volume rm $(DOCKER_CONFIG_VOLUME) || echo "No config volume to remove."
# @docker volume rm $(DOCKER_LOG_VOLUME) || echo "No log volume to remove."