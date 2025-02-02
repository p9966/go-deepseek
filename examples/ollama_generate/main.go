package main

import (
	"context"
	"fmt"
	"log"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient("")
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
}
