package httpx

import (
	"encoding/json"
	"net/http"
)

// Response is a struct for standardizing API responses
type Response struct {
	Status  int         `json:"status"`
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

// SendResponse formats and sends a JSON response with a flexible content type
func sendResponse(w http.ResponseWriter, status int, message string, data interface{}, contentType string) {
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendSuccessResponse sends a success response with a 200 status code and optional custom message
func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string, contentType ...string) {

	var (
		contentT string
	)

	if message == "" {
		message = "Success"
	}

	if len(contentType) == 0 {
		contentT = "application/json"
	} else {
		contentT = contentType[0]
	}

	sendResponse(w, http.StatusOK, message, data, contentT)
}

// SendErrorResponse sends an error response with a predefined or custom message
func SendErrorResponse(w http.ResponseWriter, code int, data interface{}, contentType ...string) {

	var (
		contentT string
	)

	if code == 0 {
		code = http.StatusInternalServerError // Default to Internal Server Error if no code is provided
	}

	if len(contentType) == 0 {
		contentT = "application/json"
	} else {
		contentT = contentType[0]
	}

	message, exists := errorMessages[code]
	if !exists {
		if code >= 400 && code < 500 {
			message = "Client Error"
		} else {
			message = "Server Error"
		}
	}
	sendResponse(w, code, message, data, contentT)
}

// CustomJSONResponse sends a custom JSON response with a specified status code
func CustomJSONResponse(w http.ResponseWriter, status int, data interface{}, contentType ...string) {
	var (
		contentT string
	)

	if len(contentType) == 0 {
		contentT = "application/json"
	} else {
		contentT = contentType[0]
	}

	w.Header().Set("Content-Type", contentT)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
