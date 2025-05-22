package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"Taurus/pkg/validate"
)

// ValidationError 用于HTTP响应的验证错误
type ValidationErrorResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Errors  map[string]string      `json:"errors,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ValidationMiddleware 创建一个HTTP请求验证中间件
// reqStruct参数是一个指向结构体的指针，用于接收和验证请求体
func ValidationMiddleware(reqStruct interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(reqStruct)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("reqStruct必须是指向结构体的指针")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 创建一个新的结构体实例
			req := reflect.New(t.Elem()).Interface()

			// 从请求体解析JSON
			if r.Body != nil && r.ContentLength > 0 {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "无法读取请求体", nil, nil)
					return
				}
				defer r.Body.Close()

				// 解析JSON到结构体
				if err := json.Unmarshal(body, req); err != nil {
					respondWithError(w, http.StatusBadRequest, "无效的JSON格式", nil, nil)
					return
				}
			}

			// 验证结构体
			if err := validate.ValidateStruct(req); err != nil {
				if valErrs, ok := err.(validate.ValidationErrors); ok {
					// 获取字段错误映射
					fieldErrors := validate.GetFieldErrors(valErrs)
					respondWithError(w, http.StatusBadRequest, "请求参数验证失败", fieldErrors, nil)
					return
				}
				// 其他错误
				respondWithError(w, http.StatusInternalServerError, "请求验证出现内部错误", nil, nil)
				return
			}

			// 将验证后的结构体存储到请求上下文
			ctx := r.Context()
			ctx = context.WithValue(ctx, "validated_request", req)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetValidatedRequest 从请求上下文中获取验证后的请求结构体
func GetValidatedRequest(r *http.Request, reqStruct interface{}) bool {
	validated := r.Context().Value("validated_request")
	if validated == nil {
		return false
	}

	// 将上下文中的值复制到提供的结构体中
	srcVal := reflect.ValueOf(validated)
	dstVal := reflect.ValueOf(reqStruct)

	if dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Struct {
		return false
	}

	dstVal.Elem().Set(srcVal.Elem())
	return true
}

// respondWithError 发送带有错误信息的HTTP响应
func respondWithError(w http.ResponseWriter, statusCode int, message string, errors map[string]string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ValidationErrorResponse{
		Code:    statusCode,
		Message: message,
		Errors:  errors,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}
