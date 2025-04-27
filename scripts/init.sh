#!/bin/bash

# 检查参数
if [ -z "$1" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

# 创建项目目录结构
mkdir -p $1/{cmd,internal,pkg,config,docs,scripts}

# 创建基本的 main.go 文件
cat <<EOL > $1/cmd/main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "your_project/pkg/router" // 替换为实际的包路径
)

func main() {
    // 初始化路由
    r := router.LoadRoutes()

    // 启动服务器
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
EOL

# 创建 README.md
echo "# $1" > $1/README.md

# 创建基本的配置文件模板
cat <<EOL > $1/config/config.yaml
# 配置文件模板
server:
  port: 8080
EOL

# 提示用户
echo "项目 $1 已初始化。" 