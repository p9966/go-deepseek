package examples

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/p9966/go-deepseek"
)

func Chat() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	request := deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, &request)
	if err != nil {
		log.Fatalf("ChatCompletion failed: %v", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
