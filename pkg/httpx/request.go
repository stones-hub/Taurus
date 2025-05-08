package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetParam 获取 GET 提交的URL参数 或 POST 提交的表单(application/x-www-form-urlencoded)数据，兼容数组
func GetParams(r *http.Request, key string) ([]string, error) {
	// 解析查询参数
	if values, ok := r.URL.Query()[key]; ok {
		return values, nil
	}

	defer r.Body.Close()

	// 解析表单数据
	if err := r.ParseForm(); err == nil {
		if values, ok := r.Form[key]; ok {
			return values, nil
		}
	}

	return nil, fmt.Errorf("key %s not found", key)
}

// GetParam 获取 GET 提交的URL参数 或 POST 提交的表单(application/x-www-form-urlencoded)数据，不兼容数组
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

// ParseJson 获取非表单提交的 JSON 对象数据, 返回map
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

// ParseJsonArray 获取非表单提交的  JSON 数组数据, 返回数组
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

// ParseJsonFlexible 获取非表单提交的  JSON 数据, 返回目标类型(不局限于数组还是对象)
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

// ParseText 获取非表单提交的纯文本数据
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

// ParseMultipartFile 解析(multipart/form-data)表单上传的文件
func ParseMultipartFile(r *http.Request, key string) ([]*multipart.FileHeader, error) {
	// 解析 multipart/form-data, 10MB 内存缓冲， 如果文件不上传完， 会报错， 所以当前函数只要返回没有错误， 就可以返回数据给客户端，不用等待
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		return nil, fmt.Errorf("failed to parse multipart form data: %w", err)
	}

	defer r.Body.Close()

	// 获取文件数据
	if files, ok := r.MultipartForm.File[key]; ok {
		return files, nil
	}

	return nil, fmt.Errorf("file key %s not found", key)
}

// ParseMultipartData 解析 multipart/form-data 请求，获取所有文件和参数数据
func ParseMultipartData(r *http.Request) (map[string][]*multipart.FileHeader, map[string][]string, error) {
	// 解析 multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		return nil, nil, fmt.Errorf("failed to parse multipart form data: %w", err)
	}
	defer r.Body.Close()

	// 获取所有文件数据
	files := r.MultipartForm.File

	// 获取参数数据
	params := r.MultipartForm.Value

	return files, params, nil
}

// SaveUploadFiles 将文件数据存储到指定目录
func SaveUploadFiles(files []*multipart.FileHeader, destDir string) error {
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		// 创建目标文件
		destPath := filepath.Join(destDir, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}
		defer destFile.Close()

		// 将上传的文件内容复制到目标文件
		if _, err := io.Copy(destFile, file); err != nil {
			return fmt.Errorf("failed to save file %s: %w", destPath, err)
		}
	}
	return nil
}

/*
- 表单post提交 并上传文件
curl -X POST \
     -H "Authorization: Bearer your_token_here" \
	 -H "Content-Type: multipart/form-data" \
     -F "file=@/path/to/file.txt" \
	 -F "file=@/path/to/file2.txt" \
	 -F "file=@/path/to/file3.txt" \
	 -F "username=testuser&password=mypassword" \
     http://example.com/api/upload

- 表单post提交 不包含文件
curl -X POST \
	 -H "Authorization: Bearer your_token_here" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "username=testuser&password=mypassword" \
     http://example.com/api/login

- json post提交
curl -X POST \
	-H "Authorization: Bearer your_token_here" \
	-H "Content-Type: application/json" \
	-d '{"username":"testuser","password":"mypassword"}' \
	http://example.com/api/login

- get请求
curl -X GET \
	-H "Authorization: Bearer your_token_here" \
    "http://example.com/api/data?user=testuser&status=active"
*/
