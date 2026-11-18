package model

import "fmt"

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

// Error Response
type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    []ValidationField `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

func (e *AppError) WithFields(fields []ValidationField) *AppError {
	e.Fields = fields
	return e
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
