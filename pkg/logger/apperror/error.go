package apperror

import(
	"fmt"
	"net/http"
)

// AppError represents an application error with context
type AppError struct {
	Err error // OG Error
	Message string // User sent message
	StatusCode int // HTTP code
	Code string // Machine-readable error code
}

func (e *AppError) Error() string {
	if e.Err != nil{
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Constructors for common errors
func NotFound(message string) *AppError{
	return &AppError{
		Message: message,
		StatusCode: http.StatusNotFound,
		Code: "NOT_FOUND",
	}
}

func BadRequest(message string, err error) *AppError{
	return &AppError{
		Message: message,
		StatusCode: http.StatusBadRequest,
		Code: "BAD_REQUEST",
	}
}

func Internal(message string, err error) *AppError{
	return &AppError{
		Message: message,
		StatusCode: http.StatusInternalServerError,
		Code: "INTERNAL_ERROR",
	}
}