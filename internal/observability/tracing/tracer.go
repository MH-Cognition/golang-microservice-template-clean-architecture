package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(getEndpoint()),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(getEnv()),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		// ✅ Sampling (100% in dev, change later if needed)
		sdktrace.WithSampler(sdktrace.AlwaysSample()),

		// ✅ Batch spans before export
		sdktrace.WithBatcher(exporter),

		// ✅ Service metadata
		sdktrace.WithResource(res),
	)

	// Set global provider
	otel.SetTracerProvider(tp)

	// ✅ Propagation (CRITICAL for microservices)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}

func getEndpoint() string {
	if ep := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); ep != "" {
		return ep
	}
	return "localhost:4317"
}

func getEnv() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	return "dev"
}
