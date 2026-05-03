package http

import (
	stderrs "errors"
	"net/http"

	"log/slog"

	appErrors "github.com/OnDev-InDev/micr_service_auth/internal/errors"
)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *appErrors.AppError

	if !stderrs.As(err, &appErr) {
		appErr = &appErrors.AppError{
			Code:    appErrors.CodeInternal,
			Message: "internal server error",
		}
	}

	logger := LoggerFromContext(r.Context())

	logger.Error("request failed",
		slog.Any("error", err),
		slog.String("code", string(appErr.Code)),
	)

	status := appErr.Code.HTTPStatus()

	jsonResponse(w, status, map[string]any{
		"error": map[string]string{
			"code":    string(appErr.Code),
			"message": appErr.Message,
		},
	})
}
