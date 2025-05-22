package middleware

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"Taurus/pkg/contextx"
	"Taurus/pkg/httpx"
	"Taurus/pkg/validate"
)

// ValidationMiddleware 创建一个HTTP请求验证中间件
func ValidationMiddleware(reqStruct interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(reqStruct)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("reqStruct必须是指向结构体的指针")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 收集所有请求数据到一个map
			data := make(map[string]interface{})

			// 1.1 收集URL查询参数
			for k, v := range r.URL.Query() {
				if len(v) == 1 {
					data[k] = v[0]
				} else {
					data[k] = v
				}
			}

			// 1.2 收集请求体数据
			if r.Body != nil && r.ContentLength > 0 {
				contentType := r.Header.Get("Content-Type")
				if idx := strings.Index(contentType, ";"); idx != -1 {
					contentType = contentType[:idx]
				}

				switch contentType {
				case "application/json":
					var jsonData map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&jsonData); err == nil {
						for k, v := range jsonData {
							data[k] = v
						}
					}
					// 恢复body以供后续读取
					r.Body = io.NopCloser(strings.NewReader(string(mustMarshal(jsonData))))

				case "application/xml":
					var xmlData map[string]interface{}
					if err := xml.NewDecoder(r.Body).Decode(&xmlData); err == nil {
						for k, v := range xmlData {
							data[k] = v
						}
					}
					r.Body = io.NopCloser(strings.NewReader(string(mustMarshal(xmlData))))

				case "application/x-www-form-urlencoded":
					if err := r.ParseForm(); err == nil {
						for k, v := range r.PostForm {
							if len(v) == 1 {
								data[k] = v[0]
							} else {
								data[k] = v
							}
						}
					}

				case "multipart/form-data":
					if err := r.ParseMultipartForm(32 << 20); err == nil {
						// 处理普通表单字段
						for k, v := range r.PostForm {
							if len(v) == 1 {
								data[k] = v[0]
							} else {
								data[k] = v
							}
						}
						// 处理文件
						if r.MultipartForm != nil {
							for k, v := range r.MultipartForm.File {
								if len(v) == 1 {
									data[k] = v[0]
								} else {
									data[k] = v
								}
							}
						}
					}
				}
			}

			// 2. 将map转换为struct
			req := reflect.New(t.Elem()).Interface()
			mapToStruct(data, req)

			// 3. 验证结构体
			if err := validate.ValidateStruct(req); err != nil {
				if valErrs, ok := err.(validate.ValidationErrors); ok {
					fieldErrors := validate.GetFieldErrors(valErrs)
					httpx.SendResponse(w, http.StatusBadRequest, fieldErrors, nil)
					return
				}
				httpx.SendResponse(w, http.StatusInternalServerError, err.Error(), nil)
				return
			}

			// 4. 将验证后的结构体存储到请求上下文
			ctx := r.Context()
			ctx = contextx.WithValidateRequest(ctx, req)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// mapToStruct 将map中的数据填充到struct
func mapToStruct(data map[string]interface{}, out interface{}) {
	v := reflect.ValueOf(out).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		// 处理json tag中的omitempty等选项
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		if val, ok := data[jsonTag]; ok {
			setFieldValueUniversal(v.Field(i), val)
		}
	}
}

// setFieldValueUniversal 支持多类型的赋值
func setFieldValueUniversal(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.String:
		switch v := value.(type) {
		case string:
			field.SetString(v)
		case []byte:
			field.SetString(string(v))
		case float64:
			field.SetString(strconv.FormatFloat(v, 'f', -1, 64))
		case int, int64:
			field.SetString(fmt.Sprintf("%v", v))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case string:
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				field.SetInt(intVal)
			}
		case float64:
			field.SetInt(int64(v))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := value.(type) {
		case string:
			if uintVal, err := strconv.ParseUint(v, 10, 64); err == nil {
				field.SetUint(uintVal)
			}
		case float64:
			field.SetUint(uint64(v))
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case string:
			if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
				field.SetFloat(floatVal)
			}
		case float64:
			field.SetFloat(v)
		}
	case reflect.Bool:
		switch v := value.(type) {
		case string:
			if boolVal, err := strconv.ParseBool(v); err == nil {
				field.SetBool(boolVal)
			}
		case bool:
			field.SetBool(v)
		}
	}
}

// mustMarshal 工具函数，将map marshal为json字符串
func mustMarshal(m map[string]interface{}) []byte {
	b, _ := json.Marshal(m)
	return b
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
