package utility

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// SuccessResponse untuk response sukses
func SuccessResponse(status int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse digunakan untuk response error
func ErrorResponse(status int, message string, errors []string) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}
