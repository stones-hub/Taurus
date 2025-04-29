#### 脚手架使用指南

---

### 一、Wire 使用

**步骤：**

1. 进入 `internal` 目录。
2. 运行 `wire` 命令。

```shell
cd internal
wire
```

---

### 二、Makefile 使用

#### 2.1、本地部署

- **启动项目**：

  ```shell
  make local-run
  ```

- **停止项目**：

  ```shell
  make local-stop
  ```

#### 2.2、本地部署（Docker 环境）

- **启动项目**：

  ```shell
  make docker-run
  ```

- **停止项目**：

  ```shell
  make docker-stop
  ```

---

### 三、Docker-Compose 部署

- **初次部署**：

  ```shell
  docker-compose up -d
  ```

- **清空容器，镜像，重新打包**：

  ```shell
  docker-compose down
  ```

- **启动项目**：

  ```shell
  docker-compose start
  ```

- **停止项目**：

  ```shell
  docker-compose stop
  ```

---

### 四、配置文件指南

- **config 目录**：用于存储应用内的各种组件的配置。添加新配置后，请在 `config/config.go` 中做好映射。

- **.env.local**：用于本地部署的默认环境变量，解决 Docker 和非 Docker 环境下参数隔离的问题。

- **.env.docker-compose**：Docker-Compose 部署所需的环境变量。使用时需复制为 `.env`。

---

### 五、初始化和更新项目

- **更新脚本**：项目更新使用的脚本是 `scripts/init.sh`。项目是否更新取决于 `.releaserc` 文件中的项目版本。

- **执行权限**：在执行 `init.sh` 之前，请确保该脚本具有执行权限。可以使用以下命令赋予权限：

  ```shell
  chmod +x scripts/init.sh
  ```

- **运行更新**：执行更新脚本以初始化或更新项目：

  ```shell
  ./scripts/init.sh
  ```

---

### 六、注意事项

- **优先级**：环境变量中的配置会覆盖 `config` 文件内的配置。

- **自定义配置路径**：可在环境变量文件中修改配置文件目录，例如：

  ```shell
  APP_CONFIG=/your_path/your_config_path
  ```

- **环境变量文件**：建议通过环境变量文件（`ENV_FILE`）传入配置，而非命令行参数。例如：

  ```shell
  make local-run ENV_FILE=/your_path/your_env_file
  ```

- **集中管理**：建议将配置集中在 `config` 目录和环境变量中进行管理。