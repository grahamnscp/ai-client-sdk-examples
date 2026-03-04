package chat

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Variables
var suseObsHTTPEndpoint = os.Getenv("SUSEOBS_EXPORTER_OTLP_HOSTNAME")
var suseObsAPIKey = os.Getenv("SUSEOBS_CLIENT_API_KEY")

// initTraceProvider create a new globally. registered trace provider
func initTraceProvider() (*trace.TracerProvider, error) {

	// Create a new OTLP trace exporter
	ctx := context.Background()

	// auth headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "SUSEObservability " + suseObsAPIKey,
	}

	// http endpoint
	// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@v1.38.0#Option
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(suseObsHTTPEndpoint),
		otlptracehttp.WithHeaders(headers),
	)
	if err != nil {
		return nil, fmt.Errorf("initTraceProvider: failed to create OTLP exporter: %w", err)
	}

	// Create a new tracer provider with the OTLP exporter.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-vllm-client"),
			semconv.ServiceVersion("1.0"),
			semconv.ServiceNamespace("local"),
			semconv.DeploymentEnvironment("dev"),
			semconv.TelemetrySDKName("openlit"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
