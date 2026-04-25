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
	CodeInternal           Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

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

func (c Code) HTTPStatus() int {
	switch c {

	case CodeInvalidCredentials,
		CodeSessionNotFound,
		CodeSessionExpired:
		return http.StatusUnauthorized

	case CodeValidationError:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
