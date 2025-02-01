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
	request := deepseek.FINCompletionRequest{
		Model:  deepseek.DeepSeekChat,
		Prompt: "What is the weather like today?",
	}

	ctx := context.Background()
	resp, err := client.CreateFINCompletion(ctx, &request)
	if err != nil {
		log.Fatalf("Error creating completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("No response choices available")
	}

	fmt.Println(resp.Choices[0].Text)
}
