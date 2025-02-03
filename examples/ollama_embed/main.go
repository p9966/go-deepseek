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
	request := deepseek.OllamaEmbedRequest{
		Model: "deepseek-r1:7b",
		Input: "Why is the sky blue?", // []string{"Why is the sky blue?", "Why is the grass green?"}
	}
	response, err := client.CreateOllamaEmbed(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create ollama embed: %v", err)
	}

	fmt.Println(response.Embeddings)
}
