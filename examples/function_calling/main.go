package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "How's the weather in Hangzhou?",
			},
		},
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

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, &request)
	if err != nil {
		log.Fatalf("ChatCompletion failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("No response choices available")
	}

	if len(resp.Choices[0].Message.ToolCalls) == 0 {
		log.Fatal("No function calls available")
	}

	fmt.Printf("Function name: %v, args:%s\n", resp.Choices[0].Message.ToolCalls[0].Function.Name, resp.Choices[0].Message.ToolCalls[0].Function.Arguments)
}
