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

	openai "github.com/sashabaranov/go-openai"
)

var OpenAITokenFile = os.Getenv("OPENAI_TOKEN_FILE")

func catchctrlc() {
	fmt.Println("\nGoodbye! If you have more questions in the future, feel free to ask. Have a great day!")
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		catchctrlc()
		os.Exit(0)
	}()

	// read openai token
	tokenFile, err := os.ReadFile(OpenAITokenFile)
	if err != nil {
		fmt.Println("failed loading openai token value: %w", err)
		return
	}
	token := string(tokenFile)
	token = strings.TrimSuffix(token, "\n")
	//fmt.Print(token)

	// new openai client connection
	client := openai.NewClient(token)

	//Model: openai.GPT3Dot5Turbo,
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4Dot1Nano,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "you are a helpful chatbot",
			},
		},
	}

	fmt.Println("\nLocal Conversation")
	fmt.Println("------------------")
	fmt.Print("> ")

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})

		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		fmt.Print("> ")
	}
}
