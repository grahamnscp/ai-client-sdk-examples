package main

import (
	"fmt"

	vllmChat "vllm-go/chat"
	vllmModels "vllm-go/models"
)

// Main
func main() {

  model := vllmModels.GetAIModel()
  fmt.Printf("vLLM AI Model: %s\n", model)
  fmt.Printf("\n")

  role := "**You are a helpful pirate chatbot. Your responses must be specific and restricted to a single line.**"

  // prompt 1
  prompt := "How deep is the pacific ocean?"
  fmt.Printf("message  >>> %s\n", prompt)
  aiResp := vllmChat.AIChat(model, role, prompt)
  fmt.Printf("response >>> %s\n", aiResp)
  fmt.Printf("\n")

  // prompt 2
  prompt = "how deep is the Mariana Trench?"
  fmt.Printf("message  >>> %s\n", prompt)
  aiResp = vllmChat.AIChat(model, role, prompt)
  fmt.Printf("response >>> %s\n", aiResp)
  fmt.Printf("\n")

  // prompt 3
  prompt = "what is the nearest city to the Mariana Trench?"
  fmt.Printf("message  >>> %s\n", prompt)
  aiResp = vllmChat.AIChat(model, role, prompt)
  fmt.Printf("response >>> %s\n", aiResp)
  fmt.Printf("\n")

}
