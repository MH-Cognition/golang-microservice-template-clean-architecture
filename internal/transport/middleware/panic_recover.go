package middleware

import (
	"fmt"
	"net/http"

	apperrors "github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/apperror"
)

// PanicRecover recovers from panics and routes them through centralized
// error handling for logging + tracing consistency.
func PanicRecover(errHandler *ErrorHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					appErr := apperrors.Wrap(
						apperrors.Internal,
						"internal server error",
						fmt.Errorf("panic: %v", rec),
					)

					if errHandler != nil {
						errHandler.handleError(w, r, appErr)
						return
					}

					// Fallback: if no handler is provided, return minimal response.
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"code":"INTERNAL_ERROR","message":"internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
