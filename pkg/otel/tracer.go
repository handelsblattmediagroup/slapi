package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func NewTracer(lc fx.Lifecycle) (trace.TracerProvider, error) {
	exporter, err := getTraceExporter()

	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if exporter != nil {
				return exporter.Shutdown(ctx)
			}
			return nil
		},
	})
	return trace.TracerProvider(tp), nil
}

func getTraceExporter() (sdktrace.SpanExporter, error) {
	tracer := os.Getenv("EOF_OTLP_TRACER")
	switch tracer {
	case "grpc":
		return otlptracegrpc.New(context.Background())
	case "http":
		return otlptracehttp.New(context.Background())
	case "zipkin":
		return zipkin.New("")
	default:
		return nil, nil
	}
}
