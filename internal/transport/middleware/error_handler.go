package middleware

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/apperror"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/logger"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HandlerFunc is a custom HTTP handler that returns error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler centralizes error handling + logging
type ErrorHandler struct {
	logger logger.Logger
}

func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}
		h.handleError(w, r, err)
	}
}

func (h *ErrorHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := apperrors.From(err)
	status := apperrors.HTTPStatus(appErr.Code)

	// ------------------------------------------------
	// OpenTelemetry: record error on span
	// ------------------------------------------------
	if span := trace.SpanFromContext(r.Context()); span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, appErr.Message)
	}

	tid := traceID(r.Context())
	sid := spanID(r.Context())

	// ------------------------------------------------
	// ERROR → root cause (DEV visibility)
	// ------------------------------------------------
	h.logger.Error(
		r.Context(),
		appErr.Error(),
		"trace_id", tid,
		"span_id", sid,
	)

	// ------------------------------------------------
	// WARN → request failed (PROD signal)
	// ------------------------------------------------
	h.logger.Warn(
		r.Context(),
		"request failed",
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"trace_id", tid,
		"span_id", sid,
	)

	// Write standardized error response (trace_id and span_id only in logs, not in response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    string(appErr.Code),
		"message": appErr.Message,
	})
}
