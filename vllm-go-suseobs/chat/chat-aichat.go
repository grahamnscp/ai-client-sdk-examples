package chat

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// variables
var apiBaseURL = fmt.Sprintf("https://%s/v1", os.Getenv("OPENAI_HOSTNAME"))
var apiKey = os.Getenv("OPENAI_API_KEY")
var tracer = otel.Tracer("vllm-client-tracer")

func AIChat(model string, role string, message string) string {

	//fmt.Printf("AIChat: Called:\n  model: %s\n  role: %s\n  message: %s\n", model, role, message)

	// Set up OpenTelemetry trace provider
	tp, err := initTraceProvider()
	if err != nil {
		log.Fatalf("AIChat: failed to initialize trace provider: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("AIChat: error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("vllm-client-tracer")
	ctx, sessionSpan := tracer.Start(context.Background(),
		"vllm-client-session",
		trace.WithAttributes(
			attribute.String("ai.model", model),
			attribute.String("user.input", message),
			attribute.String("ai.request.role", role),
			attribute.Int("ai.request.message.length", len(message)),
			attribute.String("telemetry.sdk.name", "openlit"),
		),
	)
	defer sessionSpan.End()

	// Create openapi client with config for custom baseurl and selfsigned certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	insecureClient := &http.Client{
		Transport: tr,
		Timeout:   120 * time.Second,
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = apiBaseURL
	config.HTTPClient = insecureClient

	// new openai client with config
	client := openai.NewClientWithConfig(config)

	// call local vLLM AI Server Chat Completion API..
	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: role,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	}
	// using ctx created from otel span
	resp, err := client.CreateChatCompletion(ctx, req)

	// Process response
	if err != nil {
		log.Fatalf("AIChat: ChatCompletion error: %v", err)
		sessionSpan.SetStatus(codes.Error, "ChatCompletion error")
		sessionSpan.RecordError(err)
	}

	if len(resp.Choices) > 0 {
		if resp.Usage.TotalTokens > 0 {
			sessionSpan.SetAttributes(
				attribute.String("assistant.response", resp.Choices[0].Message.Content),
				attribute.Int("ai.usage.prompt_tokens", resp.Usage.PromptTokens),
				attribute.Int("ai.usage.completion_tokens", resp.Usage.CompletionTokens),
				attribute.Int("ai.usage.total_tokens", resp.Usage.TotalTokens),
			)
		}
		sessionSpan.SetStatus(codes.Ok, "Success")
		return resp.Choices[0].Message.Content
	}
	sessionSpan.SetStatus(codes.Error, "No response received")
	return "No response received."
}
