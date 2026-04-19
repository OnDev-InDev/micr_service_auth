package http_layer

import (
	"errors"
	"net/http"

	"micr_service_auth/internal/service"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func formatValidationErrors(errs validation.Errors) map[string]string {
	res := make(map[string]string)

	for field, err := range errs {
		if err != nil {
			res[field] = err.Error()
		}
	}

	return res
}

func writeError(w http.ResponseWriter, err error) {
	var ve validation.Errors

	switch {
	case errors.As(err, &ve):
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"errors": formatValidationErrors(ve),
		})

	case errors.Is(err, service.ErrInvalidCredentials):
		jsonResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid email or password",
		})

	default:
		jsonResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
}
