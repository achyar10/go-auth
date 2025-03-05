package utility

import "time"

type APIResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    []string    `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// SuccessResponse untuk response sukses
func SuccessResponse(status int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// ErrorResponse digunakan untuk response error
func ErrorResponse(status int, message string, errors []string) APIResponse {
	return APIResponse{
		Status:    status,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
