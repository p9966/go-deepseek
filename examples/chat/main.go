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
		Model: deepseek.DeepseekV4Pro,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		ExtraBody: deepseek.ExtraBody{
			Thinking: deepseek.ExtraBodyThinking{
				Type: deepseek.ExtraBodyThinkingModelEnabled,
			},
		},
		ReasoningEffort: deepseek.ReasoningEffortModelMax,
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, &request)
	if err != nil {
		log.Fatalf("ChatCompletion failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("No response choices available")
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
