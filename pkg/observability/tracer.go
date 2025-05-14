package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/michelsazevedo/hawkeye/internal/config"
)

func NewTracerProvider(conf *config.Config) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(conf.Observability.JaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracer := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(conf.Settings.ApplicationName),
		)),
	)

	otel.SetTracerProvider(tracer)
	return tracer, nil
}
