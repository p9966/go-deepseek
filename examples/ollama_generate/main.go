package main

import (
	"context"
	"fmt"
	"log"

	"github.com/p9966/go-deepseek"
)

const (
	baseURL = "http://localhost:11434"
	model   = "deepseek-r1:7b"
)

func main() {
	client := deepseek.Client{
		BaseUrl: baseURL,
	}

	// 示例 1：简单生成请求
	executeGenerateRequest(client, createGenerateRequest("Why is the sky blue?", nil))

	// 示例 2：结构化输出
	structuredFormat := map[string]interface{}{
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
	}
	executeGenerateRequest(client, createGenerateRequest("simple is 29 years old and is busy saving the world. Respond using JSON", structuredFormat))

	// 示例 3：JSON 模式
	executeGenerateRequest(client, createGenerateRequest("What color is the sky at different times of the day? Respond using JSON", "json"))
}

// createGenerateRequest 根据给定参数创建一个新的 OllamaGenerateRequest。
func createGenerateRequest(prompt string, format any) deepseek.OllamaGenerateRequest {
	return deepseek.OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: format,
	}
}

// executeGenerateRequest 执行生成请求并打印响应结果。
func executeGenerateRequest(client deepseek.Client, request deepseek.OllamaGenerateRequest) {
	resp, err := client.CreateOllamaGenerate(context.TODO(), &request)
	if err != nil {
		log.Fatalf("failed to create generate: %v", err)
	}
	fmt.Println(resp.Response)
}
