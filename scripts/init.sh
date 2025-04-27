#!/bin/bash

# 定义颜色
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
BLUE='\033[34m'
CYAN='\033[36m'
RESET='\033[0m'

# 显示帮助信息
show_help() {
  echo -e "${CYAN}欢迎使用 Taurus 框架管理脚本${RESET}"
  echo -e "${CYAN}========================================${RESET}"
  echo -e "${GREEN}可用命令：${RESET}"
  echo -e "${YELLOW}install${RESET} - 安装框架到指定目录"
  echo -e "${YELLOW}update${RESET}  - 更新框架到最新版本或指定版本"
  echo -e "${CYAN}========================================${RESET}"
  echo -e "${GREEN}用法示例：${RESET}"
  echo -e "${YELLOW}./init.sh install${RESET} - 安装框架"
  echo -e "${YELLOW}./init.sh update${RESET}  - 更新框架"
  echo -e "${CYAN}========================================${RESET}"
}

# 检查是否提供了命令
if [ -z "$1" ]; then
  show_help
  exit 1
fi

# 函数：检查输入是否为空
check_empty_input() {
  local input="$1"
  local prompt="$2"
  while [ -z "$input" ]; do
    echo "输入不能为空，请重新输入。"
    read -p "$prompt" input
  done
  echo "$input"
}

# 函数：获取最新发布版本
get_latest_release() {
  curl --silent "https://api.github.com/repos/stones-hub/Taurus/releases/latest" | # 获取最新发布版本信息
    grep '"tag_name":' |                                                          # 提取 tag_name 字段
    sed -E 's/.*"([^"]+)".*/\1/'                                                  # 提取版本号
}

# 函数：获取当前版本
get_current_version() {
  local project_path="$1"
  # 从指定路径的 .releaserc 文件中获取版本号，如果获取不到则使用默认版本号 v1.0.0
  VERSION=$(grep -Eo 'version:[^ ]+' "$project_path/.releaserc" | sed -E 's/version:(.*)/\1/' || echo "v1.0.0")
  echo "$VERSION"
}

# 函数：比较版本号
version_greater_equal() {
  # 使用 sort -V 比较版本号
  [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

# 函数：更新配置文件中的指定项
update_config_items() {
  local version="$1"
  local authorization="$2"
  local db_enable="$3"
  local redis_enable="$4"
  local print_config="$5"

  # 使用完整的项目路径
  local config_path="$project_path/$project_name/config"

  # 遍历 config 目录中的 yaml, json, toml 文件
  find "$config_path" -type f \( -name "*.yaml" -o -name "*.json" -o -name "*.toml" \) | while read -r file; do
    case "$file" in
      *.yaml)
        sed -i.bak -E "s/(version: *\").*\"/\1$version\"/" "$file"
        sed -i.bak -E "s/(authorization: *\").*\"/\1$authorization\"/" "$file"
        sed -i.bak -E "s/(db_enable: *).*/\1$db_enable/" "$file"
        sed -i.bak -E "s/(redis_enable: *).*/\1$redis_enable/" "$file"
        sed -i.bak -E "s/(print_config: *).*/\1$print_config/" "$file"
        ;;
      *.json)
        sed -i.bak -E "s/(\"version\": *\").*\"/\1$version\"/" "$file"
        sed -i.bak -E "s/(\"authorization\": *\").*\"/\1$authorization\"/" "$file"
        sed -i.bak -E "s/(\"db_enable\": *).*/\1$db_enable,/" "$file"
        sed -i.bak -E "s/(\"redis_enable\": *).*/\1$redis_enable,/" "$file"
        sed -i.bak -E "s/(\"print_config\": *).*/\1$print_config/" "$file"
        ;;
      *.toml)
        sed -i.bak -E "s/(version = *\").*\"/\1$version\"/" "$file"
        sed -i.bak -E "s/(authorization = *\").*\"/\1$authorization\"/" "$file"
        sed -i.bak -E "s/(db_enable = *).*/\1$db_enable/" "$file"
        sed -i.bak -E "s/(redis_enable = *).*/\1$redis_enable/" "$file"
        sed -i.bak -E "s/(print_config = *).*/\1$print_config/" "$file"
        ;;
    esac
    rm "${file}.bak" # 删除备份文件
  done
}

# 函数：安装框架
install_framework() {
  # 定义颜色
  RED='\033[31m'
  GREEN='\033[32m'
  YELLOW='\033[33m'
  BLUE='\033[34m'
  CYAN='\033[36m'
  RESET='\033[0m'

  # 第一步：设置项目路径
  echo -e "${CYAN}请输入项目下载路径（默认: 当前目录）:${RESET}"
  read -p "" project_path
  project_path=${project_path:-.}
  project_path=$(check_empty_input "$project_path" "请输入项目下载路径（默认: 当前目录）: ")

  # 展开路径中的 ~ 并去掉末尾的斜杠
  project_path=$(eval echo "$project_path" | sed 's:/*$::')

  # 检查目录是否存在
  if [ ! -d "$project_path" ]; then
    echo -e "${YELLOW}目录 $project_path 不存在，是否创建? (y/n):${RESET}"
    read -p "" create_dir
    if [ "$create_dir" == "y" ]; then
      mkdir -p "$project_path" || { echo -e "${RED}无法创建目录，程序退出。${RESET}"; exit 1; }
    else
      echo -e "${RED}目录不存在且未创建，程序退出。${RESET}"
      exit 1
    fi
  fi

  # 第二步：设置项目名称
  echo -e "${CYAN}请输入项目名称（默认: Taurus_demo）:${RESET}"
  read -p "" project_name
  project_name=${project_name:-Taurus_demo}

  # 处理项目名称中的特殊字符
  project_name=$(echo "$project_name" | sed 's/[ /\\]/_/g')
  project_name=$(check_empty_input "$project_name" "项目名称不能为空，请重新输入: ")

  # 第三步：选择框架版本
  echo -e "${CYAN}请输入框架版本（默认: latest）:${RESET}"
  read -p "" framework_version
  framework_version=${framework_version:-latest}

  # 检查版本是否存在
  if [ "$framework_version" != "latest" ]; then
    available_versions=$(curl --silent "https://api.github.com/repos/stones-hub/Taurus/tags" | grep '"name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if ! echo "$available_versions" | grep -q "^$framework_version$"; then
      echo -e "${YELLOW}版本号 $framework_version 不存在，正在下载最新版本...${RESET}"
      framework_version="latest"
    fi
  fi

  # 从 GitHub 下载框架模板
  echo -e "${BLUE}正在从 GitHub 下载框架模板...${RESET}"
  if [ "$framework_version" == "latest" ]; then
    curl -L -o "$project_path/$project_name.tar.gz" "https://github.com/stones-hub/Taurus/archive/refs/heads/main.tar.gz"
  else
    curl -L -o "$project_path/$project_name.tar.gz" "https://github.com/stones-hub/Taurus/archive/refs/tags/$framework_version.tar.gz"
  fi

  # 解压下载的文件
  mkdir -p "$project_path/$project_name" || { echo -e "${RED}无法创建项目目录，程序退出。${RESET}"; exit 1; }
  tar -xzf "$project_path/$project_name.tar.gz" -C "$project_path/$project_name" --strip-components=1 || { echo -e "${RED}解压失败，程序退出。${RESET}"; exit 1; }
  rm "$project_path/$project_name.tar.gz"

  # 第四步：设置授权码
  echo -e "${CYAN}请输入授权码（默认: Bearer 123456）:${RESET}"
  read -p "" authorization
  authorization=${authorization:-"Bearer 123456"}

  # 第五步：询问是否启用数据库
  echo -e "${CYAN}是否启用数据库? (y/n):${RESET}"
  read -p "" use_db
  if [ "$use_db" == "y" ]; then
    db_enable=true
  else
    db_enable=false
  fi

  # 第六步：询问是否启用 Redis
  echo -e "${CYAN}是否启用 Redis? (y/n):${RESET}"
  read -p "" use_redis
  if [ "$use_redis" == "y" ]; then
    redis_enable=true
  else
    redis_enable=false
  fi

  # 第七步：是否打印配置
  echo -e "${CYAN}是否打印配置? (y/n):${RESET}"
  read -p "" print_config
  if [ "$print_config" == "y" ]; then
    print_config=true
  else
    print_config=false
  fi

  # 更新配置文件中的指定项
  update_config_items "$framework_version" "$authorization" "$db_enable" "$redis_enable" "$print_config"

  echo -e "${GREEN}项目已成功初始化在 $project_path/$project_name 目录下。${RESET}"
}

# 函数：更新框架
update_framework() {
  # 定义颜色
  RED='\033[31m'
  GREEN='\033[32m'
  YELLOW='\033[33m'
  BLUE='\033[34m'
  CYAN='\033[36m'
  RESET='\033[0m'

  # 定义分割符
  SEPARATOR="${CYAN}========================================${RESET}"

  # 提示用户输入本地框架根目录
  echo -e "$SEPARATOR"
  echo -e "${CYAN}请输入本地框架的根目录：${RESET}"
  read -p "" project_path

  # 展开路径中的 ~ 并去掉末尾的斜杠
  project_path=$(eval echo "$project_path" | sed 's:/*$::')

  # 检查是否提供了本地框架根目录
  if [ -z "$project_path" ]; then
    echo -e "${RED}错误：未提供本地框架的根目录，程序退出。${RESET}"
    exit 1
  fi

  # 检查项目是否存在
  echo -e "$SEPARATOR"
  if [ ! -f "$project_path/Makefile" ] || [ ! -d "$project_path/cmd" ]; then
    echo -e "${RED}错误：项目不存在，程序退出。${RESET}"
    exit 1
  fi

  # 备份项目目录
  echo -e "$SEPARATOR"
  echo -e "${BLUE}正在备份项目目录...${RESET}"
  backup_file="${project_path}_backup_$(date +%Y%m%d%H%M%S).tar.gz"
  tar -czf "$backup_file" -C "$project_path" . || { echo -e "${RED}备份失败，程序退出。${RESET}"; exit 1; }
  echo -e "${GREEN}项目已备份到 $backup_file${RESET}"

  # 提示用户输入版本号
  echo -e "$SEPARATOR"
  echo -e "${CYAN}请输入要更新到的版本号（或按 Enter 键以更新到最新版本）：${RESET}"
  read -p "" version

  # 如果用户未输入版本号，获取最新发布版本
  if [ -z "$version" ]; then
    echo -e "${YELLOW}未输入版本号，正在获取最新发布版本...${RESET}"
    version=$(get_latest_release)
    if [ -z "$version" ]; then
      echo -e "${RED}无法获取最新发布版本，程序退出。${RESET}"
      exit 1
    fi
    echo -e "${GREEN}最新发布版本为 $version${RESET}"
  fi

  # 获取当前版本
  echo -e "$SEPARATOR"
  current_version=$(get_current_version "$project_path")
  echo -e "${BLUE}当前版本为 $current_version${RESET}"

  # 检查版本号
  if version_greater_equal "$current_version" "$version"; then
    echo -e "${YELLOW}当前版本已是最新版本或高于目标版本，无法更新。${RESET}"
    exit 1
  fi

  # 下载新的框架版本
  echo -e "$SEPARATOR"
  echo -e "${BLUE}正在下载框架版本 $version...${RESET}"
  curl -L -o "update.tar.gz" "https://github.com/stones-hub/Taurus/archive/refs/tags/$version.tar.gz" || { echo -e "${RED}下载失败，程序退出。${RESET}"; exit 1; }

  # 解压到临时目录
  echo -e "$SEPARATOR"
  mkdir -p "update_temp" || { echo -e "${RED}无法创建临时目录，程序退出。${RESET}"; exit 1; }
  tar -xzf "update.tar.gz" -C "update_temp" --strip-components=1 || { echo -e "${RED}解压失败，程序退出。${RESET}"; exit 1; }
  rm "update.tar.gz"

  # 强制更新 pkg 、 script、internal/app 目录
  echo -e "$SEPARATOR"
  echo -e "${BLUE}更新 pkg 、 script、internal/app 目录...${RESET}"
  rsync -aq update_temp/pkg "$project_path/pkg"
  rsync -aq update_temp/scripts "$project_path/scripts"
  rsync -aq update_temp/internal/app "$project_path/internal/app"

  # 强制更新 internal 目录下的 injector.go 和 wire.go 文件
  echo -e "$SEPARATOR"
  echo -e "${BLUE}更新 internal 目录下的 injector.go 和 wire.go 文件...${RESET}"
  rsync -aq update_temp/internal/injector.go "$project_path/internal/injector.go"
  rsync -aq update_temp/internal/wire.go "$project_path/internal/wire.go"

  # 强制更新 config 目录下的 config.go 文件
  echo -e "$SEPARATOR"
  echo -e "${BLUE}更新 config 目录下的 config.go 文件...${RESET}"
  rsync -aq update_temp/config/config.go "$project_path/config/config.go"

  # 更新根目录下的文件
  echo -e "$SEPARATOR"
  echo -e "${BLUE}更新根目录下的文件... ${RESET}"
  rsync -aq update_temp/.dockerignore "$project_path/.dockerignore"
  rsync -aq update_temp/.gitignore "$project_path/.gitignore"
  rsync -aq update_temp/.releaserc "$project_path/.releaserc"
  rsync -aq update_temp/Dockerfile "$project_path/Dockerfile"
  rsync -aq update_temp/go.mod "$project_path/go.mod"
  rsync -aq update_temp/go.sum "$project_path/go.sum"
  rsync -aq update_temp/LICENSE "$project_path/LICENSE"
  rsync -aq update_temp/Makefile "$project_path/Makefile"
  rsync -aq update_temp/README.md "$project_path/README.md"
  
  # 清理临时文件
  echo -e "$SEPARATOR"
  rm -rf "update_temp"

  echo -e "${GREEN}框架已更新到版本: $version ${RESET}"
  echo -e "$SEPARATOR"
}

# 主程序逻辑
case "$1" in
  install)
    install_framework
    ;;
  update)
    update_framework
    ;;
  *)
    show_help
    exit 1
    ;;
esac 