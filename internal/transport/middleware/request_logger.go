package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// traceID extracts trace_id from context
func traceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// spanID extracts span_id from context
func spanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

// statusRecorder captures HTTP status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RequestLogger logs request lifecycle information
func RequestLogger(log logger.Logger, env string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			ctx := r.Context()
			tid := traceID(ctx)
			sid := spanID(ctx)

			// ------------------------------------------------
			// INFO: Request received
			// ------------------------------------------------
			log.Info(ctx, r.Method+" "+r.URL.Path, "trace_id", tid, "span_id", sid)

			// ------------------------------------------------
			// DEBUG: Raw request body (DEV only)
			// ------------------------------------------------
			if env != "prod" && r.Body != nil {
				body, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				if len(body) > 0 {
					log.Debug(ctx, "Raw request body received", "trace_id", tid, "span_id", sid)
				}
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// ------------------------------------------------
			// Span enrichment
			// ------------------------------------------------
			span := trace.SpanFromContext(ctx)
			if span.SpanContext().IsValid() {
				span.SetAttributes(
					attribute.Int("http.status_code", rec.status),
					attribute.Int64("http.duration_ms", duration.Milliseconds()),
				)
			}

			// ------------------------------------------------
			// INFO: Request completed (DEV visibility)
			// ------------------------------------------------
			log.Info(ctx, "request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration.Milliseconds(),
				"trace_id", tid,
				"span_id", sid,
			)
		})
	}
}

