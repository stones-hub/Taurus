# 默认环境变量文件
ENV_FILE=.env.local

# 检查环境变量文件是否存在
ifeq ($(wildcard $(ENV_FILE)),)
$(error $(RED)Environment file '$(ENV_FILE)' not found. Please create it or specify a different file.$(RESET))
endif

# ---------- 注意事项 ----------
# 1. 凡是通过make命令运行的程序，都会加载ENV_FILE中的环境变量，所以程序启动的时候，环境变量已经有了
# 2. make local-run 还是传递了-env参数，只不过是为了兼容应用程序而已，其实可以不传
# 3. make docker-run 传递了 --env-file $(ENV_FILE) 是因为docker容器中的环境变量跟宿主机不互通，所以需要传，但是我在Dockerfile中并没有让容器运行的时候传递-env参数， 是因为--env-file $(ENV_FILE) 会自动将ENV_FILE中的环境变量写入到容器中
# 4. 用法 # make run-local  ENV_FILE=.env.local,  make docker-run ENV_FILE=.env.local 指定环境变量文件  
#
# -----------------------------

# ---------------------------- 加载环境变量 --------------------------------
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

# 定义颜色
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
CYAN := \033[36m
RESET := \033[0m
# 定义分割符
SEPARATOR := $(CYAN)--------------------------------$(RESET)

# ---------------------------- 从环境变量中获取配置 --------------------------------

# 检查环境变量是否为空
ifeq ($(strip $(VERSION)),)
$(error $(RED)VERSION is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(APP_NAME)),)
$(error $(RED)APP_NAME is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(APP_HOST)),)
$(error $(RED)APP_HOST is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(APP_PORT)),)
$(error $(RED)APP_PORT is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(APP_CONFIG)),)
$(error $(RED)APP_CONFIG is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(HOST_PORT)),)
$(error $(RED)HOST_PORT is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(CONTAINER_PORT)),)
$(error $(RED)CONTAINER_PORT is not set in the environment variables$(RESET))
endif

ifeq ($(strip $(WORKDIR)),)
$(error $(RED)WORKDIR is not set in the environment variables$(RESET))
endif

BUILD_DIR := build
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_CONTAINER := $(APP_NAME)
DOCKER_NETWORK := $(APP_NAME)-network
DOCKER_LOG_VOLUME := $(APP_NAME)_log_data

# ---------------------------- 构建目标 --------------------------------
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
local-run: clean build
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Running the application locally...$(RESET)"
	@$(BUILD_DIR)/$(APP_NAME) -config=$(APP_CONFIG) -env=$(ENV_FILE) || echo -e "$(RED)Failed to run the application locally.$(RESET)"
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
	@echo -e "$(BLUE)Checking if Docker log volume exists...$(RESET)"
	@if ! docker volume inspect $(DOCKER_LOG_VOLUME) >/dev/null 2>&1; then \
		echo -e "$(YELLOW)Creating Docker volume for logs...$(RESET)"; \
		docker volume create $(DOCKER_LOG_VOLUME); \
	fi
	@echo -e "$(BLUE)Running the application in Docker...$(RESET)"
	@docker run -d --name $(DOCKER_CONTAINER) --network $(DOCKER_NETWORK) -p ${HOST_PORT}:${CONTAINER_PORT} \
		--env-file $(ENV_FILE) \
		--mount type=bind,source=$(shell realpath $(APP_CONFIG)),target=$(WORKDIR)/config \
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

# docker-compose.yml
docker-compose-up:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Starting Docker Compose...$(RESET)"
	@docker-compose -f docker-compose.yml up -d || echo -e "$(RED)Failed to start Docker Compose.$(RESET)"
	@echo -e "$(GREEN)Docker Compose started.$(RESET)"
	@echo -e "$(SEPARATOR)"

# docker-compose.yml
docker-compose-down:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Stopping Docker Compose...$(RESET)"
	@docker-compose -f docker-compose.yml down || echo -e "$(RED)Failed to stop Docker Compose.$(RESET)"
	@echo -e "$(GREEN)Docker Compose stopped.$(RESET)"
	@echo -e "$(SEPARATOR)"

# docker-compose-swarm.yml
docker-compose-swarm-up: docker-build
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Pushing Docker image to registry...$(RESET)"
	echo docker tag $(DOCKER_IMAGE) $(REGISTRY_URL)/$(DOCKER_IMAGE)
	echo docker push $(REGISTRY_URL)/$(DOCKER_IMAGE) || echo -e "$(RED)Failed to push Docker image.$(RESET)"
	@echo -e "$(GREEN)Docker image pushed to registry.$(RESET)"
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Starting Docker Compose Swarm...$(RESET)"
	echo docker-compose -f docker-compose-swarm.yml up -d || echo -e "$(RED)Failed to start Docker Compose Swarm.$(RESET)"
	@echo -e "$(GREEN)Docker Compose Swarm started.$(RESET)"
	@echo -e "$(SEPARATOR)"

# docker-compose-swarm.yml
docker-compose-swarm-down:
	@echo -e "$(SEPARATOR)"
	@echo -e "$(BLUE)Stopping Docker Compose Swarm...$(RESET)"
	@docker-compose -f docker-compose-swarm.yml down || echo -e "$(RED)Failed to stop Docker Compose Swarm.$(RESET)"
	@echo -e "$(GREEN)Docker Compose Swarm stopped.$(RESET)"