package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
  "os/signal"
	"strings"
  "syscall"

	"google.golang.org/genai"
	//"google.golang.org/api/option"
)

func catchctrlc() {
  fmt.Println("\nGoodbye! note; you can type quit or exit to stop the program.")
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

  // check token value present in default variable and populate var
  geminiApiKey := os.Getenv("GEMINI_API_KEY")
	if geminiApiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set. Please set it before running the program.")
	}

	// Initialize the Gemini client
	ctx := context.Background()

	//client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
    APIKey:   geminiApiKey,
	})
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}


	// Select Generative Model
	// "gemini-pro", "gemini-1.5-flash", "gemini-1.5-pro"

	// Configure model parameters (temperature for creativity etc)
	//model.SetTemperature(0.7) // Higher values (closer to 1.0) for more creative outputs
	//model.SetMaxOutputTokens(800) // Limit the length of responses

	// Start chat session, maintaining conversation history across turns
	//cs := model.StartChat()
  cs, _ := client.Chats.Create(ctx, "gemini-1.5-flash", nil, nil)

	fmt.Println("LocalChatbot: Hello! I'm your AI assistant. Type 'quit' to exit.")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userMessage := scanner.Text()

    switch strings.ToLower(userMessage) {
		  case "quit","exit": {
			  fmt.Println("Goodbye!")
			  os.Exit(0)
      }
      case "": break
      default: {
     
		    // Send message and get a response
		    resp, err := cs.SendMessage(ctx, genai.Part{Text: userMessage})
		    if err != nil {
			    log.Printf("Error sending message to Gemini: %v", err)
			    fmt.Println("LocalChatbot: Encountered an error. Please try again.")
			    fmt.Print("> ")
			    continue
		    }

		    // Process and print the response
        if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			    // Iterate through each "Part" in the content of the first candidate.
			    // The content might contain text, images, or other types of data.
			    for _, part := range resp.Candidates[0].Content.Parts {
            if part.Text != "" { // If the 'Text' field is not empty, it's a text part.
					    fmt.Printf("LocalChatbot:\n%s\n", part.Text)
				    }
			    }
		    } else {
			    fmt.Println("LocalChatbot: Couldn't generate a response. Please try rephrasing.")
		    }
      } //end case default
    } // switch stdin

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

