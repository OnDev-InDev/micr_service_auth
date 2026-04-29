package http

import (
	"net/http"
  stderrs "errors"

	appErrors "micr_service_auth/internal/errors"
)

func writeError(w http.ResponseWriter, err error) {
	var appErr *appErrors.AppError

	if !stderrs.As(err, &appErr) {
		appErr = &appErrors.AppError{
			Code:    appErrors.CodeInternal,
			Message: "internal server error",
		}
	}

	status := appErr.Code.HTTPStatus()

	jsonResponse(w, status, map[string]any{
		"error": map[string]string{
			"code":    string(appErr.Code),
			"message": appErr.Message,
		},
	})
}