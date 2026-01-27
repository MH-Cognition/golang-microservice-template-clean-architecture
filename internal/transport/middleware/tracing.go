package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func Tracing(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.Tracer(serviceName)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := otel.GetTextMapPropagator().Extract(
				r.Context(),
				propagation.HeaderCarrier(r.Header),
			)

			ctx, span := tracer.Start(
				ctx,
				r.Method+" "+r.URL.Path,
			)
			defer span.End()

			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.path", r.URL.Path),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}