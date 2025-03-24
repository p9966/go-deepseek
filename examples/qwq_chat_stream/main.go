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
	client := deepseek.Client{
		BaseUrl:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
		AuthToken: os.Getenv("QwQ_AUTH_TOKEN"), // 获取地址：https://bailian.console.aliyun.com/?apiKey=1#/api-key
	}
	scanner := bufio.NewScanner(os.Stdin)

	var messages []deepseek.ChatCompletionMessage

	//  To do console input in debug mode, add "console": "integratedTerminal" to launch.json
	for {

		fmt.Print("You: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatalf("Failed to read input: %v", err)
			}
			break
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
			Model:    deepseek.QwQ_plus_latest,
			Messages: messages,
		}

		ctx := context.Background()
		stream, err := client.CreateChatCompletionStream(ctx, request)
		if err != nil {
			log.Fatalf("ChatCompletionStream failed: %v", err)
		}
		defer stream.Close()

		fmt.Print("QWQ: ")
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println()
				break
			}
			if err != nil {
				log.Fatalf("ChatCompletionStream stream.Recv() failed: %v", err)
			}

			if len(response.Choices) > 0 && response.Choices[0].FinishReason == "" {
				content := response.Choices[0].Delta.Content
				messages = append(messages, deepseek.ChatCompletionMessage{
					Role:    deepseek.ChatMessageRoleAssistant,
					Content: content,
				})
				fmt.Print(content)
			}
		}
	}
}
