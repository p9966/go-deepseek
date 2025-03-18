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
		Model: deepseek.QWen2_5_7b,
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
						Properties: map[string]any{
							"location": map[string]any{
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

	funcMap := map[string]func(...string) string{}
	funcMap["get_weather"] = getWeather
	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create ollama embed: %v", err)
	}
	if response.Message.ToolCalls != nil {
		if getWeater, ok := funcMap[response.Message.ToolCalls[0].Function.Name]; ok {
			result := getWeater(response.Message.ToolCalls[0].Function.Arguments["location"].(string))
			fmt.Println(result)
		}
	}
}

// functions
func getWeather(args ...string) string {
	location := args[0]
	return fmt.Sprintf("The weather in %s is sunny", location)
}
