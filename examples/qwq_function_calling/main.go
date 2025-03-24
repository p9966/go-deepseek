package main

import (
	"context"
	"encoding/json"
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
		AuthToken: os.Getenv("QWQ_AUTH_TOKEN"), // 获取地址：https://bailian.console.aliyun.com/?apiKey=1#/api-key
	}
	request := deepseek.StreamChatCompletionRequest{
		Model: deepseek.QwQ_32b,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "成都天气怎么样",
			},
		},
		Tools: []deepseek.Tools{
			{
				Type: "function",
				Function: deepseek.Function{
					Name:        "get_weather",
					Description: "当你想查询指定城市的天气时非常有用",
					Parameters: &deepseek.Parameters{
						Type: "object",
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"description": "城市或县区，比如北京市、杭州市、成都市、余杭区等",
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
			content := response.Choices[0].Delta.ReasoningContent

			fmt.Print(content)

			if len(response.Choices[0].Delta.ToolCalls) > 0 {
				if buf, err := json.Marshal(response.Choices[0].Delta.ToolCalls); err == nil {
					fmt.Println(string(buf))
				}
			}
		}
	}
}
