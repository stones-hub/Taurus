# gRPC 验证拦截器示例

这个示例展示了如何在 gRPC 服务中使用验证拦截器来验证请求参数。

## 目录结构

```
grpc_validator/
├── Makefile
├── proto/
│   └── user/
│       └── user.proto
├── server/
│   └── main.go
└── client/
    └── main.go
```

## 使用方法

### 1. 定义 proto 文件

在 proto 文件中定义服务接口和消息类型：

```protobuf
syntax = "proto3";

package user;
option go_package = "Taurus/example/grpc_validator/proto/user";

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse) {}
}

message CreateUserRequest {
  string name = 1;  // 用户名
  string email = 2; // 邮箱
  int32 age = 3;    // 年龄
  string password = 4; // 密码
}

message CreateUserResponse {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}

message UpdateUserRequest {
  int64 id = 1;     // 用户ID
  string name = 2;  // 用户名
  string email = 3; // 邮箱
  int32 age = 4;    // 年龄
}

message UpdateUserResponse {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}
```

### 2. 生成 Go 代码

运行以下命令生成 Go 代码：

```bash
make proto
```

### 3. 添加验证标签

在生成的 Go 代码中，为请求结构体添加验证标签：

```go
type CreateUserRequest struct {
    Name     string `validate:"required,min=2,max=50"`
    Email    string `validate:"required,email"`
    Age      int32  `validate:"required,gt=0,lt=150"`
    Password string `validate:"required,min=6,max=20"`
}

type UpdateUserRequest struct {
    Id    int64  `validate:"required,gt=0"`
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int32  `validate:"required,gt=0,lt=150"`
}
```

### 4. 在服务端使用验证拦截器

在服务端代码中注册验证拦截器：

```go
srv, cleanup, err := server.NewServer(
    server.WithAddress(":50051"),
    server.WithUnaryInterceptor(
        interceptor.UnaryServerValidationInterceptor(),
    ),
)
```

### 5. 运行示例

1. 启动服务端：
```bash
make server
```

2. 在另一个终端运行客户端：
```bash
make client
```

## 验证规则说明

验证拦截器使用 `validate` 标签来定义验证规则，常用的规则包括：

- `required`: 字段必填
- `min=n`: 最小长度/值
- `max=n`: 最大长度/值
- `gt=n`: 大于 n
- `lt=n`: 小于 n
- `email`: 邮箱格式
- `url`: URL 格式
- `numeric`: 数字格式
- `alpha`: 字母格式
- `alphanumeric`: 字母数字格式

## 错误处理

当请求参数验证失败时，服务端会返回 `InvalidArgument` 错误，错误消息包含具体的验证失败原因，例如：

```
请求参数验证失败: name长度必须大于2; email必须是有效的邮箱格式; age必须大于0; password长度必须大于6
```

## 注意事项

1. 确保在生成的 Go 代码中添加了正确的验证标签
2. 验证标签的规则要与业务需求相匹配
3. 验证失败时会直接返回错误，不会继续处理请求
4. 验证拦截器会自动处理所有请求，不需要在每个服务方法中手动验证 