package main

import (
	"context"
	"fmt"
	"log"

	"github.com/p9966/go-deepseek"
)

const (
	baseURL = "http://localhost:11434"
	model   = "deepseek-r1:7b"
)

func main() {
	client := deepseek.Client{
		BaseUrl: baseURL,
	}

	// Example 1: Simple generate request
	executeGenerateRequest(client, createGenerateRequest("Why is the sky blue?", nil))

	// Example 2: Structured outputs
	structuredFormat := map[string]interface{}{
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
	}
	executeGenerateRequest(client, createGenerateRequest("simple is 29 years old and is busy saving the world. Respond using JSON", structuredFormat))

	// Example 3: JSON mode
	executeGenerateRequest(client, createGenerateRequest("What color is the sky at different times of the day? Respond using JSON", "json"))
}

// createGenerateRequest creates a new OllamaGenerateRequest with the given parameters.
func createGenerateRequest(prompt string, format any) deepseek.OllamaGenerateRequest {
	return deepseek.OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: format,
	}
}

// executeGenerateRequest executes the generate request and prints the response.
func executeGenerateRequest(client deepseek.Client, request deepseek.OllamaGenerateRequest) {
	resp, err := client.CreateOllamaGenerate(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Response)
}
