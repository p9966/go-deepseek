package main

import (
	"context"
	"fmt"
	"log"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.Client{
		BaseUrl: "http://localhost:11434",
	}

	// Example 1: Simple chat completion
	simpleRequest := createChatRequest("deepseek-r1:7b", "why is the sky blue?", nil, nil)
	executeChatRequest(client, simpleRequest)

	// Example 2: Structured outputs
	structuredRequest := createChatRequest(
		"deepseek-r1:7b",
		"simple is 29 years old and is busy saving the world. Respond using JSON",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
				"age": map[string]interface{}{
					"type": "integer",
				},
				"available": map[string]interface{}{
					"type": "boolean",
				},
			},
			"required": []string{"name", "age", "available"},
		},
		nil,
	)
	executeChatRequest(client, structuredRequest)
}

// createChatRequest creates a new OllamaChatRequest with the given parameters.
func createChatRequest(model, content string, format map[string]interface{}, tools []deepseek.Tools) deepseek.OllamaChatRequest {
	return deepseek.OllamaChatRequest{
		Model: model,
		Messages: []deepseek.OllamaChatMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
		Format: format,
		Tools:  tools,
	}
}

// executeChatRequest executes the chat request and prints the response.
func executeChatRequest(client deepseek.Client, request deepseek.OllamaChatRequest) {
	resp, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Message.Content)
}
