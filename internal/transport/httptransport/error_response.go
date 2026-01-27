package httptransport

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/apperror"
)


// ErrorResponse is the standard HTTP error response contract
// This structure MUST remain stable across services.
// Note: trace_id and span_id are only logged, not included in HTTP responses.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteError writes a JSON error response to the client.
// It should be called ONLY from centralized error handling middleware.
func WriteError(w http.ResponseWriter, err error) {
	appErr := apperrors.From(err)
	if appErr == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status := apperrors.HTTPStatus(appErr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
	})
}