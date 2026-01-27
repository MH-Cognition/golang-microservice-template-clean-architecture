package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type otelTracer struct{}

func NewTracer() Tracer {
	return &otelTracer{}
}

func (t *otelTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	ctx, span := otel.Tracer("app").Start(ctx, name)
	return ctx, spanWrapper{span: span}
}

type spanWrapper struct {
	span trace.Span
}

func (s spanWrapper) End() {
	s.span.End()
}
