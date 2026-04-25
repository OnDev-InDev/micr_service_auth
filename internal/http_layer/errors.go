package http_layer

import (
	"net/http"

	"micr_service_auth/internal/errors"
)

func writeError(w http.ResponseWriter, err error) {

	appErr, ok := err.(*errors.AppError)
	if !ok {
		jsonResponse(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    string(errors.CodeInternal),
				"message": "internal server error",
			},
		})
		return
	}

	status := appErr.Code.HTTPStatus()

	jsonResponse(w, status, map[string]any{
		"error": map[string]string{
			"code":    string(appErr.Code),
			"message": appErr.Message,
		},
	})
}
