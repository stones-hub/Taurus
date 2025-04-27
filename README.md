
## 脚手架使用指南


### 一、Wire使用

1. 进入internal目录，使用wire命令


### 二、Makefile使用

#### 名词解释
- HOST_PROT 
    主机映射端口
- CONTAINER_PORT 
    容器内应用端口
- WORKDIR
    容器内工作目录

#### 项目启动方式

1. 容器管理

- 启动
    > HOST_PORT=8090 CONTAINER_PORT=8090 WORKDIR=/myapp make docker-run 
- 停止
    > make docker-stop

2. 本地部署

- 启动
    > make local-run

- 停止
    > make local-stop

- 清理
    > make clean