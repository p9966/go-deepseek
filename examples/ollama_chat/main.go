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

	// 示例 1：简单聊天补全
	simpleRequest := createChatRequest("deepseek-r1:7b", "why is the sky blue?", nil, nil)
	executeChatRequest(client, simpleRequest)

	// 示例 2：结构化输出
	structuredRequest := createChatRequest(
		"deepseek-r1:7b",
		"simple is 29 years old and is busy saving the world. Respond using JSON",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
				"age": map[string]interface{}{
					"type": "integer",
				},
				"available": map[string]interface{}{
					"type": "boolean",
				},
			},
			"required": []string{"name", "age", "available"},
		},
		nil,
	)
	executeChatRequest(client, structuredRequest)
}

// createChatRequest 根据给定参数创建一个新的 OllamaChatRequest。
func createChatRequest(model, content string, format map[string]interface{}, tools []deepseek.Tools) deepseek.OllamaChatRequest {
	return deepseek.OllamaChatRequest{
		Model: model,
		Messages: []deepseek.OllamaChatMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
		Format: format,
		Tools:  tools,
	}
}

// executeChatRequest 执行聊天请求并打印响应结果。
func executeChatRequest(client deepseek.Client, request deepseek.OllamaChatRequest) {
	resp, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Message.Content)
}
