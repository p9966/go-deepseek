package deepseek

import (
	"context"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPSEEK_API_KEY not set")
	}

	client := NewClient(apiKey)
	request := ChatCompletionRequest{
		Model: DeepSeekChat,
		Messages: []ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices returned")
	}

	t.Logf("Response: %v", resp.Choices[0].Message.Content)
}
