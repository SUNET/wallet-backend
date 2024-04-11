package trace

import (
	"context"
	"time"
	"vc/pkg/logger"
	"vc/pkg/model"

	jaegerPropagator "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is a wrapper for opentelemetry tracer
type Tracer struct {
	TP *sdktrace.TracerProvider
	trace.Tracer
	log *logger.Log
}

func newExporter(ctx context.Context, cfg *model.Cfg) (sdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.Common.Tracing.Addr),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(10*time.Second),
	)
}

func newTraceProvider(exp sdktrace.SpanExporter, serviceName string) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
}

// New return a new tracer
func New(ctx context.Context, cfg *model.Cfg, log *logger.Log, projectName, serviceName string) (*Tracer, error) {
	exp, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tracer := &Tracer{
		TP:  newTraceProvider(exp, serviceName),
		log: log,
	}

	otel.SetTracerProvider(tracer.TP)
	otel.SetTextMapPropagator(jaegerPropagator.Jaeger{})

	tracer.Tracer = otel.Tracer(projectName)

	return tracer, nil
}

// Shutdown shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	t.log.Info("Shutting down tracer")
	if err := t.TP.Shutdown(ctx); err != nil {
		return err
	}
	ctx.Done()
	return nil
}
