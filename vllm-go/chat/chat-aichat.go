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
)

// variables
var apiBaseURL = fmt.Sprintf("https://%s/v1", os.Getenv("OPENAI_HOSTNAME"))
var apiKey = os.Getenv("OPENAI_API_KEY")

func AIChat(model string, role string, message string) string {

	//fmt.Printf("AIChat: Called:\n  model: %s\n  role: %s\n  message: %s\n", model, role, message)

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
	resp, err := client.CreateChatCompletion(context.Background(), req)

  // Process response
	if err != nil {
		log.Fatalf("AIChat: ChatCompletion error: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return "No response received."
}
