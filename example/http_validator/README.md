# HTTP 验证中间件示例

这个示例展示了如何在 HTTP 服务中使用验证中间件来验证请求参数。

## 目录结构

```
http_validator/
├── controller/
│   └── user.go
├── middleware/
│   └── validator.go
├── main.go
└── README.md
```

## 功能特点

1. 支持多种请求格式：
   - JSON
   - XML
   - Form 表单
   - URL 查询参数
   - 文件上传

2. 自动收集和合并请求数据：
   - 从请求体
   - 从 URL 查询参数
   - 从表单数据
   - 从文件上传

3. 灵活的验证规则：
   - 必填验证
   - 长度验证
   - 数值范围验证
   - 格式验证（邮箱、URL等）
   - 自定义验证规则

## 使用方法

### 1. 定义请求结构体

在控制器中定义请求结构体，添加验证标签：

```go
type UserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int32  `json:"age" validate:"required,gt=0,lt=150"`
    Password string `json:"password" validate:"required,min=6,max=20"`
}
```

### 2. 注册路由和中间件

在路由注册时使用验证中间件：

```go
mux.Handle("/api/user/create", 
    middleware.ValidationMiddleware(&UserRequest{})(
        http.HandlerFunc(userCtrl.CreateUser),
    ),
)
```

### 3. 在控制器中获取验证后的数据

```go
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    req, ok := contextx.GetValidateRequest(r.Context()).(*UserRequest)
    if !ok {
        httpx.SendResponse(w, http.StatusBadRequest, "无效的请求数据", nil)
        return
    }
    // 使用验证后的数据...
}
```

## 验证规则说明

验证中间件使用 `validate` 标签来定义验证规则，常用的规则包括：

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

## 示例请求

### 创建用户

```bash
curl -X POST http://localhost:8080/api/user/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "age": 25,
    "password": "123456"
  }'
```

### 更新用户

```bash
curl -X POST http://localhost:8080/api/user/update \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com",
    "age": 26
  }'
```

## 错误处理

当请求参数验证失败时，服务端会返回 400 Bad Request 状态码，响应体包含具体的验证失败原因，例如：

```json
{
  "code": 400,
  "message": {
    "name": "name长度必须大于2",
    "email": "email必须是有效的邮箱格式",
    "age": "age必须大于0",
    "password": "password长度必须大于6"
  }
}
```

## 注意事项

1. 确保请求结构体中的字段标签（json、validate）正确设置
2. 验证规则要与业务需求相匹配
3. 验证失败时会直接返回错误，不会继续处理请求
4. 验证中间件会自动处理所有请求数据，不需要手动解析和验证
5. 支持多种请求格式的混合使用，如同时使用 JSON 和查询参数 