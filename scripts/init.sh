#!/bin/bash

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
  # 从 .releaserc 文件中获取版本号，如果获取不到则使用默认版本号 v1.0.0
  VERSION=$(grep -Eo 'version:[^ ]+' .releaserc | sed -E 's/version:(.*)/\1/' || echo "v1.0.0")
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
  local version="$1"
  if [ -z "$version" ]; then
    echo "未输入版本号，正在获取最新发布版本..."
    version=$(get_latest_release)
    if [ -z "$version" ]; then
      echo "无法获取最新发布版本，程序退出。"
      exit 1
    fi
    echo "最新发布版本为 $version"
  fi

  # 获取当前版本
  current_version=$(get_current_version)
  echo "当前版本为 $current_version"

  # 检查版本号
  if version_greater_equal "$current_version" "$version"; then
    echo "当前版本已是最新版本或高于目标版本，无法更新。"
    exit 1
  fi

  # 下载新的框架版本
  echo "正在下载框架版本 $version..."
  curl -L -o "update.tar.gz" "https://github.com/stones-hub/Taurus/archive/refs/tags/$version.tar.gz" || { echo "下载失败，程序退出。"; exit 1; }

  # 解压到临时目录
  mkdir -p "update_temp" || { echo "无法创建临时目录，程序退出。"; exit 1; }
  tar -xzf "update.tar.gz" -C "update_temp" --strip-components=1 || { echo "解压失败，程序退出。"; exit 1; }
  rm "update.tar.gz"

  # 覆盖更新
  echo "正在更新框架..."
  rsync -av --exclude='cmd' --exclude='config' --exclude='internal' --exclude='logs' update_temp/ "$project_path/"
  rsync -av update_temp/internal/app/ "$project_path/internal/app/"

  # 更新配置文件中的版本号
  update_config_items "$version" "$authorization" "$db_enable" "$redis_enable" "$print_config"

  # 清理临时文件
  rm -rf "update_temp"

  echo "框架已更新到版本 $version。"
}

# 主程序逻辑
case "$1" in
  install)
    install_framework
    ;;
  update)
    update_framework "$2"
    ;;
  *)
    echo "用法: $0 {install|update <version>}"
    exit 1
    ;;
esac 