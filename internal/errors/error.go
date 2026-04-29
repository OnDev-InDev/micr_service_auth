package errors

import (
	"fmt"
	"net/http"
)

type Code string

const (
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	CodeSessionNotFound    Code = "SESSION_NOT_FOUND"
	CodeSessionExpired     Code = "SESSION_EXPIRED"
	CodeValidationError    Code = "VALIDATION_ERROR"
	CodeUnauthorized       Code = "UNAUTHORIZED"
	CodeForbidden          Code = "FORBIDDEN"
)

const (
	ErrDB    Code = "DB_ERROR"
	ErrJSON  Code = "JSON_ERROR"
)

const (
	CodeMethodNotAllowed Code = "METHOD_NOT_ALLOWED"
	CodeBadRequest       Code = "BAD_REQUEST"
	CodeInternal         Code = "INTERNAL SERVER"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}


// метод реализующий интерфейс 
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code Code, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

// маппинг
func (c Code) HTTPStatus() int {
	switch c {

	case CodeInvalidCredentials,
		CodeSessionNotFound,
		CodeSessionExpired,
		CodeUnauthorized:
		return http.StatusUnauthorized

	case CodeForbidden:
		return http.StatusForbidden

	case CodeValidationError:
		return http.StatusBadRequest

	case CodeMethodNotAllowed:
		return http.StatusMethodNotAllowed

	default:
		return http.StatusInternalServerError
	}
}