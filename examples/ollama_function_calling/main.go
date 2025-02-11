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
	request := deepseek.OllamaChatRequest{
		Model: "llama3.1:8b",
		Messages: []deepseek.OllamaChatMessage{
			{
				Role:    "user",
				Content: "How's the weather in chengdu?",
			},
		},
		Stream: false,
		Tools: []deepseek.Tools{
			{
				Type: "function",
				Function: deepseek.Function{
					Name:        "get_weather",
					Description: "Get weather of an location, the user shoud supply a location first",
					Parameters: &deepseek.Parameters{
						Type: "object",
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"description": "The location to get weather",
								"type":        "string",
							},
						},
						Required: []string{"location"},
					},
				},
			},
		},
	}
	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create ollama embed: %v", err)
	}
	fmt.Println(response.Message.ToolCalls[0].Function.Name, response.Message.ToolCalls[0].Function.Arguments)
}
