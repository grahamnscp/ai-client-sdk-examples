package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var apiURL = fmt.Sprintf("http://%s/ollama/v1/chat/completions", 
                             os.Getenv("OPEN_WEBUI_HOSTNAME"))

var apiKey = os.Getenv("OPEN_WEBUI_API_KEY")

// Message represents a single message in the chat, with a role and content
// The role can be "user", "assistant", or "system"
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the JSON payload sent to the OpenWebUI API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents the JSON response received from the API
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Main
func main() {
	// Initialize a list to hold the chat history
	var chatHistory []Message

	// Add a system message to set the context for the model
	chatHistory = append(chatHistory, Message{
		Role:    "system",
		Content: "You are a helpful assistant.",
	})

	fmt.Println("Chat with OpenWebUI. Type 'quit' or 'exit' to end the session.")
	reader := bufio.NewReader(os.Stdin)

	// Create a custom HTTP client that skips TLS certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

  // loop on chat..
	for {
		// Prompt the user for input
		fmt.Print("> ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		userInput = userInput[:len(userInput)-1] // Remove the newline

		// Check for exit commands.
		if userInput == "quit" || userInput == "exit" {
			fmt.Println("Exiting chat.")
			break
		}

		// Append message to the chat history
		chatHistory = append(chatHistory, Message{
			Role:    "user",
			Content: userInput,
		})

		// Create the request payload with model name
		requestPayload := ChatRequest{
			Model:    "gemma:7B", 
			Messages: chatHistory,
		}

		// Marshal the request payload into a JSON byte array
		jsonPayload, err := json.Marshal(requestPayload)
		if err != nil {
			log.Fatalf("Error marshaling JSON payload: %v", err)
		}

		// Create a new HTTP POST request
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Fatalf("Error creating HTTP request: %v", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		// Send request..
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending HTTP request: %v", err)
		}
		defer resp.Body.Close()

		// Check for HTTP errors
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Fatalf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Unmarshal the JSON response into ChatResponse struct
		var responsePayload ChatResponse
		err = json.Unmarshal(body, &responsePayload)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON response: %v", err)
		}

		// Check if a response message was found
		if len(responsePayload.Choices) > 0 {
			assistantMessage := responsePayload.Choices[0].Message

			// Print the model's response.
			fmt.Printf("Response: %s\n", assistantMessage.Content)

			// Append the response message to the chat history for context
			chatHistory = append(chatHistory, assistantMessage)
		} else {
			fmt.Println("Response: No response received.")
		}
	}
}

