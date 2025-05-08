package httpx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Response is a struct for standardizing API responses
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	StatusInvalidRequest = 1001 // Invalid Request
	StatusInvalidParams  = 1002 // Invalid Parameters
	StatusUnauthorized   = 1003 // Unauthorized
)

// ErrorMessages holds custom error messages for different error codes
var errorMessages = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusNotImplemented:      "Not Implemented",
	http.StatusBadGateway:          "Bad Gateway",
	http.StatusServiceUnavailable:  "Service Unavailable",
	StatusInvalidRequest:           "Invalid Request",    // 无效请求
	StatusInvalidParams:            "Invalid Parameters", // 无效参数
	StatusUnauthorized:             "Unauthorized",       // 未授权
}

// SendResponse formats and sends a response with a flexible content type
func SendResponse(w http.ResponseWriter, code int, data interface{}, headers map[string]string) {
	httpStatus, message := getResponseStatusAndMessage(code)

	// 如果 headers 为 nil，初始化为一个空的 map
	if headers == nil {
		headers = make(map[string]string)
	}

	// 如果 headers 中没有 Content-Type，默认设置为 application/json;charset=utf-8
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json;charset=utf-8"
	}

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	contentType := headers["Content-Type"]

	// 写入响应头
	w.WriteHeader(httpStatus)

	// 根据不同类型的 contentType 前缀，进行不同的数据处理, 支持 xml/json/text/html
	if strings.HasPrefix(contentType, "application/json") {
		json.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
	} else if strings.HasPrefix(contentType, "application/xml") {
		xml.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
	} else if strings.HasPrefix(contentType, "text/plain") || strings.HasPrefix(contentType, "text/html") {
		if str, ok := data.(string); ok {
			w.Write([]byte(str))
		} else {
			// 将 data 转换为 JSON 字符串
			jsonData, err := json.Marshal(data)
			if err != nil {
				w.Write([]byte("Response Error converting data to JSON"))
			} else {
				w.Write(jsonData)
			}
		}
	} else {
		// 默认返回 json
		json.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
	}
}

// CustomJSONResponse sends a custom JSON response with a specified status code
func CustomJSONResponse(w http.ResponseWriter, data interface{}, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func getResponseStatusAndMessage(code int) (httpStatus int, message string) {
	if _, ok := errorMessages[code]; ok {
		return http.StatusOK, errorMessages[code]
	}

	if http.StatusText(code) != "" {
		return code, http.StatusText(code)
	}

	return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
}

// RedirectResponse sends a redirect response to the client
func RedirectResponse(w http.ResponseWriter, url string, code int) {
	if code < 300 || code > 399 {
		code = http.StatusFound // 默认使用 302 Found
	}
	http.Redirect(w, &http.Request{}, url, code)
}

// HTMLResponse sends an HTML response to the client
func HTMLResponse(w http.ResponseWriter, htmlContent string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlContent))
}

// FileResponseWithManualRangeSupport sends a file to the client for download with manual range support
func FileResponseWithManualRangeSupport(w http.ResponseWriter, r *http.Request, filePath string, fileName string) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not obtain file info.", http.StatusInternalServerError)
		return
	}

	// 解析 Range 请求头
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// 如果没有 Range 请求头，直接发送整个文件
		sendFullFile(w, file, fileInfo, fileName)
		return
	}

	// 处理 Range 请求头
	start, end, err := parseRange(rangeHeader, fileInfo.Size())
	if err != nil {
		http.Error(w, "Invalid Range header.", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.WriteHeader(http.StatusPartialContent)

	// 发送文件部分
	file.Seek(start, io.SeekStart)
	io.CopyN(w, file, end-start+1)
}

// sendFullFile 发送整个文件
func sendFullFile(w http.ResponseWriter, file *os.File, fileInfo os.FileInfo, fileName string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	io.Copy(w, file)
}

// parseRange 解析 Range 请求头
func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	rangeParts := strings.Split(rangeHeader, "=")
	if len(rangeParts) != 2 || rangeParts[0] != "bytes" {
		return 0, 0, fmt.Errorf("invalid range")
	}

	byteRange := strings.Split(rangeParts[1], "-")
	start, err := strconv.ParseInt(byteRange[0], 10, 64)
	if err != nil || start < 0 || start >= fileSize {
		return 0, 0, fmt.Errorf("invalid range start")
	}

	var end int64
	if byteRange[1] != "" {
		end, err = strconv.ParseInt(byteRange[1], 10, 64)
		if err != nil || end < start || end >= fileSize {
			return 0, 0, fmt.Errorf("invalid range end")
		}
	} else {
		end = fileSize - 1
	}

	return start, end, nil
}

// FileDownloadWithRange  客户端断点续传下载文件
func FileDownloadWithRange(url, destPath string) error {
	// 打开文件，准备写入
	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// 获取文件的当前大小
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}
	currentSize := fileInfo.Size()

	// 创建请求，并设置 Range 头
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentSize))

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server does not support partial content, status code: %d", resp.StatusCode)
	}

	// 将响应写入文件
	_, err = file.Seek(currentSize, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking in file: %v", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

/*

// 下载文件服务端
http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
	// 假设文件路径和文件名是通过查询参数传递的
	filePath := r.URL.Query().Get("file")
	fileName := filepath.Base(filePath) // 从路径中提取文件名

	// 使用 FileResponseWithManualRangeSupport 处理文件下载
	httpx.FileResponseWithManualRangeSupport(w, r, filePath, fileName)
})

// 下载文件客户端
func main() {
	url := "http://localhost:8080/download?file=/path/to/your/file.txt"
	destPath := filepath.Join(".", "file.txt")

	err := httpx.FileDownloadWithRange(url, destPath)
	if err != nil {
		fmt.Println("Download failed:", err)
	} else {
		fmt.Println("Download completed successfully")
	}
}

*/
