package helper

import (
	"net/http"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
)

func NewAppError(statusCode int, code, message string, err error) *model.AppError {
	return &model.AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func ErrBadRequest(code, message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusBadRequest,
		Code:       code,
		Message:    message,
	}
}

func ErrUnauthorized(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func ErrForbidden(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

func ErrNotFound(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func ErrNotAcceptable(message string) *model.AppError {
	if message == "" {
		message = "header Accept tidak didukung oleh endpoint ini"
	}
	return &model.AppError{
		StatusCode: http.StatusNotAcceptable,
		Code:       "NOT_ACCEPTABLE",
		Message:    message,
	}
}

func ErrUnsupportedMediaType(message string) *model.AppError {
	if message == "" {
		message = "Content-Type wajib application/json"
	}
	return &model.AppError{
		StatusCode: http.StatusUnsupportedMediaType,
		Code:       "UNSUPPORTED_MEDIA_TYPE",
		Message:    message,
	}
}

func ErrConflict(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func ErrValidation(message string, fields []model.ValidationField) *model.AppError {
	if message == "" {
		message = "data masukan tidak valid"
	}
	return &model.AppError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Fields:     fields,
	}
}

func ErrServiceUnavailable(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusServiceUnavailable,
		Code:       "SERVICE_UNAVAILABLE",
		Message:    message,
	}
}

func ErrRateLimited(message string) *model.AppError {
	if message == "" {
		message = "terlalu banyak percobaan, silakan coba lagi nanti"
	}
	return &model.AppError{
		StatusCode: http.StatusTooManyRequests,
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    message,
	}
}

func ErrInternal(err error) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "terjadi kesalahan internal pada server",
		Err:        err,
	}
}

func ErrTooManyRequest(message string) *model.AppError {
	return &model.AppError{
		StatusCode: http.StatusTooManyRequests,
		Code:       "TOO_MANY_REQUEST",
		Message:    message,
	}
}

func ToValidationFields(errs map[string]string) []model.ValidationField {
	vals := []model.ValidationField{}

	for field, err := range errs {
		val := model.ValidationField{
			Field:   field,
			Message: err,
		}

		vals = append(vals, val)
	}

	return vals
}
