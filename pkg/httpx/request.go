package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// GetParam 获取 GET 或 POST 提交的数据，兼容数组
func GetParams(r *http.Request, key string) ([]string, error) {
	// 解析查询参数
	if values, ok := r.URL.Query()[key]; ok {
		return values, nil
	}

	// 解析表单数据
	if err := r.ParseForm(); err == nil {
		if values, ok := r.Form[key]; ok {
			return values, nil
		}
	}

	return nil, fmt.Errorf("key %s not found", key)
}

// GetParam 获取 GET 或 POST 提交的数据，不兼容数组
func GetParam(r *http.Request, key string) (string, error) {
	if res, err := GetParams(r, key); err != nil {
		return "", err
	} else {
		if len(res) == 0 {
			return "", fmt.Errorf("key %s not found", key)
		}

		return res[0], nil
	}
}

// ParseUploadFile 解析上传的文件
func ParseUploadFile(r *http.Request, key string) ([]*multipart.FileHeader, error) {
	// 解析 multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		return nil, fmt.Errorf("failed to parse multipart form data: %w", err)
	}

	// 获取文件数据
	if files, ok := r.MultipartForm.File[key]; ok {
		return files, nil
	}

	return nil, fmt.Errorf("file key %s not found", key)
}

// ParseJson 获取提交的 JSON 数据
func ParseJson(r *http.Request) (map[string]interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return nil, fmt.Errorf("content type is not application/json")
	}

	defer r.Body.Close()
	var jsonData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %w", err)
	}

	return jsonData, nil
}

// ParseText 获取提交的纯文本数据
func ParseText(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		return "", fmt.Errorf("content type is not text/plain")
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read plain text body: %w", err)
	}

	return string(body), nil
}

// ParseJsonArray 获取提交的 JSON 数组数据
func ParseJsonArray(r *http.Request) ([]interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return nil, fmt.Errorf("content type is not application/json")
	}

	defer r.Body.Close()
	var jsonArray []interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jsonArray); err != nil {
		return nil, fmt.Errorf("failed to parse JSON array body: %w", err)
	}

	return jsonArray, nil
}

// ParseJsonFlexible 根据传入的目标类型解析 JSON 数据
func ParseJsonFlexible(r *http.Request, target interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("content type is not application/json")
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to parse JSON body: %w", err)
	}

	return nil
}
