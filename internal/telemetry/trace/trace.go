package trace

import (
	"log"

	"github.com/amirhnajafiz/hls/internal/telemetry/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// New tracer.
func New(cfg config.Trace) trace.Tracer {
	if !cfg.Enabled {
		return trace.NewNoopTracerProvider().Tracer("fhs")
	}

	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Agent.Host), jaeger.WithAgentPort(cfg.Agent.Port)),
	)
	if err != nil {
		log.Fatalf("failed to initialize export pipeline: %v", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semConv.ServiceNamespaceKey.String("amirhnajafiz"),
			semConv.ServiceNameKey.String("file-handling-system"),
		),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdkTrace.NewBatchSpanProcessor(exporter)
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.ParentBased(sdkTrace.TraceIDRatioBased(cfg.Ratio))),
		sdkTrace.WithSpanProcessor(bsp),
		sdkTrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	// register the TraceContext propagator globally.
	var tc propagation.TraceContext

	otel.SetTextMapPropagator(tc)

	tracer := otel.Tracer("fhs")

	return tracer
}
