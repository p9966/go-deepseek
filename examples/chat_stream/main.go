package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("hi:")

	var messages []deepseek.ChatCompletionMessage

	for {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		messages = append(messages, deepseek.ChatCompletionMessage{
			Role:    deepseek.ChatMessageRoleUser,
			Content: input,
		})

		request := deepseek.StreamChatCompletionRequest{
			Model:    deepseek.DeepSeekChat,
			Messages: messages,
		}

		ctx := context.Background()

		stream, err := client.CreateChatCompletionStream(ctx, request)
		if err != nil {
			log.Fatalf("ChatCompletionStream failed: %v", err)
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatalf("ChatCompletionStream stream.Recv() failed: %v", err)
			}

			if response.Choices != nil {
				if response.Choices[0].FinishReason != "" {
					break
				}
				messages = append(messages, deepseek.ChatCompletionMessage{
					Role:    deepseek.ChatMessageRoleAssistant,
					Content: response.Choices[0].Delta.Content,
				})

				fmt.Print(response.Choices[0].Delta.Content)
			}
		}
	}
}
