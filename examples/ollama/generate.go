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
	request := deepseek.OllamaGenerateRequest{
		Model:  "deepseek-r1:7b", // or other local models
		Prompt: "Why is the sky blue?",
		Stream: false,
	}

	resp, err := client.CreateOllamaGenerate(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Response)

	// Structured outputs
	request.Prompt = "simple is 29 years old and is busy saving the world. Respond using JSON"
	request.Format = map[string]interface{}{
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
	resp, err = client.CreateOllamaGenerate(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Response)

	// JSON mode
	request.Prompt = "What color is the sky at different times of the day? Respond using JSON"
	request.Format = "json"
	resp, err = client.CreateOllamaGenerate(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Response)
}
