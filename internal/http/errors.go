package http

import (
	stderrs "errors"
	"log"
	"net/http"

	appErrors "micr_service_auth/internal/errors"
)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *appErrors.AppError

	if !stderrs.As(err, &appErr) {
		// если не признаем ошибку, то маскируем под это
		appErr = &appErrors.AppError{
			Code:    appErrors.CodeInternal,
			Message: "internal server error",
		}
	}

	requestID, _ := r.Context().Value(RequestIDKey).(string)

	log.Printf(
		"request_id=%s error=%v",
		requestID,
		err,
	)

	status := appErr.Code.HTTPStatus()

	jsonResponse(w, status, map[string]any{
		"error": map[string]string{
			"code":    string(appErr.Code),
			"message": appErr.Message,
		},
	})
}
