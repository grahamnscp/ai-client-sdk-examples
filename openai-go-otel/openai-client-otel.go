// openai golang

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	//openai "github.com/sashabaranov/go-openai"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"

	otelopenai "github.com/langwatch/langwatch/sdk-go/instrumentation/openai"
)

var OpenAITokenFile = os.Getenv("OPENAI_TOKEN_FILE")
var OTelExporterOLTPEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
var OTelExporterOLTPTracesEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")

func catchctrlc() {
	fmt.Println("\nGoodbye! If you have more questions in the future, feel free to ask. Have a great day!")
}

// MAIN
func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		catchctrlc()
		os.Exit(0)
	}()

	// read openai token from file
	tokenFile, err := os.ReadFile(OpenAITokenFile)
	if err != nil {
		fmt.Println("failed loading openai token value: %w", err)
		return
	}
	token := string(tokenFile)
	token = strings.TrimSuffix(token, "\n")


	// Setup tracer provider
	ctx := context.Background()
	setupOTelTraceProvider(ctx)

	// openai client with langwatch tracing
	client := openai.NewClient(
		option.WithAPIKey(token),
		option.WithMiddleware(otelopenai.Middleware("localchatapp",
			// Optional: Capture request/response content (be mindful of sensitive data)
			otelopenai.WithCaptureInput(),
			otelopenai.WithCaptureOutput(),
		)),
	)

	// conversation system prompt
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful chatbot."),
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4_1Nano,
	}

	req, err := client.Chat.Completions.New(ctx, param)
	if err != nil {
		fmt.Printf("Chat completion failed: %v", err)
		return
	}

	// loop on standard in
	fmt.Println("\nLocal Conversation")
	fmt.Println("------------------")
	fmt.Print("> ")

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
    switch strings.ToLower(s.Text()) {
      case "quit","exit": {
        fmt.Println("Goodbye!")
        os.Exit(0)
      }
    }
		param.Messages = append(param.Messages, req.Choices[0].Message.ToParam())
		param.Messages = append(param.Messages, openai.UserMessage(s.Text()))

		resp, err := client.Chat.Completions.New(ctx, param)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		param.Messages = append(param.Messages, resp.Choices[0].Message.ToParam())
		fmt.Print("> ")
	}
}

func setupOTelTraceProvider(ctx context.Context) func() {

  // Set OTLP traces endpoint from env
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(OTelExporterOLTPTracesEndpoint),
		//otlptracehttp.WithHeaders(map[string]string{"Authorization": "Bearer " + ObsAPIKey}),
	)
	if err != nil {
		fmt.Printf("failed to create OTLP exporter: %v\n", err)
		os.Exit(1)
	}

	// Set default SDK resources and the required service name are set.
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("LocalGoChatApp"),
		),
	)
	if err != nil {
		fmt.Printf("failed to create OTLP SDK resource: %v\n", err)
		os.Exit(1)
	}

  // Initialise and set trace provider with exporter and resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)

	// Return a function to shutdown the tracer provider
	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown TracerProvider: %v\n", err)
			os.Exit(1)
		}
	}
}
