package apperrors

import "github.com/ayo-awe/go-backend-starter/internal/oapi"

type AppError struct {
	// We maintain the dependency on api.ErrorCode because we want it to match the opeanpi spec verbatim
	Code       oapi.ErrorCode
	StatusCode int
	Message    string
}

func (e AppError) Error() string {
	return e.Message
}

// create a new AppError instance
func New(message string, code oapi.ErrorCode, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
	}
}
