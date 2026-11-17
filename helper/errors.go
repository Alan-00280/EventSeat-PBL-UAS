package helper

import (
	"fmt"
	"net/http"
)

type ValidationField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	StatusCode int               `json:"-"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Fields     []ValidationField `json:"fields,omitempty"`
	Err        error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d %s] %s: %v", e.StatusCode, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d %s] %s", e.StatusCode, e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func (e *AppError) WithFields(fields []ValidationField) *AppError {
	e.Fields = fields
	return e
}

func ErrBadRequest(code, message string) *AppError {
	if code == "" {
		code = "BAD_REQUEST"
	}
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       code,
		Message:    message,
	}
}

func ErrUnauthorized(code, message string) *AppError {
	if code == "" {
		code = "UNAUTHORIZED"
	}
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       code,
		Message:    message,
	}
}

func ErrForbidden(code, message string) *AppError {
	if code == "" {
		code = "FORBIDDEN"
	}
	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       code,
		Message:    message,
	}
}

func ErrNotFound(code, message string) *AppError {
	if code == "" {
		code = "NOT_FOUND"
	}
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       code,
		Message:    message,
	}
}

func ErrNotAcceptable(message string) *AppError {
	if message == "" {
		message = "header Accept tidak didukung oleh endpoint ini"
	}
	return &AppError{
		StatusCode: http.StatusNotAcceptable,
		Code:       "NOT_ACCEPTABLE",
		Message:    message,
	}
}

func ErrUnsupportedMediaType(message string) *AppError {
	if message == "" {
		message = "Content-Type wajib application/json"
	}
	return &AppError{
		StatusCode: http.StatusUnsupportedMediaType,
		Code:       "UNSUPPORTED_MEDIA_TYPE",
		Message:    message,
	}
}

func ErrConflict(code, message string) *AppError {
	if code == "" {
		code = "CONFLICT"
	}
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       code,
		Message:    message,
	}
}

func ErrValidation(message string, fields []ValidationField) *AppError {
	if message == "" {
		message = "data masukan tidak valid"
	}
	return &AppError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Fields:     fields,
	}
}

func ErrRateLimited(message string) *AppError {
	if message == "" {
		message = "terlalu banyak percobaan, silakan coba lagi nanti"
	}
	return &AppError{
		StatusCode: http.StatusTooManyRequests,
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    message,
	}
}

func ErrInternal(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "terjadi kesalahan internal pada server",
		Err:        err,
	}
}

func ErrTooManyRequest(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusTooManyRequests,
		Code:       "TOO_MANY_REQUEST",
		Message:    message,
	}
}
