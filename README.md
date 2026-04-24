# Go DeepSeek
[![Go Reference](https://pkg.go.dev/badge/github.com/p9966/go-deepseek.svg)](https://pkg.go.dev/github.com/p9966/go-deepseek)
[![Go Report Card](https://goreportcard.com/badge/github.com/p9966/go-deepseek)](https://goreportcard.com/report/github.com/p9966/go-deepseek)

`go-deepseek` is an unofficial Go client for DeepSeek-compatible APIs. It can be used with DeepSeek, Qwen3/QwQ via DashScope-compatible endpoints, OpenAI-compatible APIs, and local Ollama models.

## Features
- Chat completions
- Streaming chat completions
- Function calling / tool calling
- FIM (Fill-in-Middle) completions
- Embeddings through Ollama
- DeepSeek account balance query
- Configurable `BaseUrl`, timeout, and custom `http.Client`

## Installation
```bash
go get github.com/p9966/go-deepseek
```

Requires Go `1.23+`.

## Quick Start
This is the fastest way to call the default DeepSeek API:

```go
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

	req := deepseek.ChatCompletionRequest{
		Model: deepseek.DeepseekV4Pro,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "Explain why the sky is blue in one paragraph.",
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("no response choices returned")
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

## Configure Another Compatible Endpoint
The client defaults to `https://api.deepseek.com`, but you can point it to any compatible endpoint.

### Qwen3 / QwQ via DashScope
```go
client := deepseek.NewClient(
	os.Getenv("QWEN3_AUTH_TOKEN"),
	deepseek.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
)
```

### OpenAI-compatible endpoint
```go
client := deepseek.NewClient(
	os.Getenv("OPENAI_API_KEY"),
	deepseek.WithBaseURL("https://api.openai.com/v1"),
)
```

### Custom timeout or HTTP client
```go
client := deepseek.NewClient(
	os.Getenv("DEEPSEEK_API_KEY"),
	deepseek.WithTimeout(30*time.Second),
)
```

Or:

```go
client := deepseek.NewClient(
	os.Getenv("DEEPSEEK_API_KEY"),
	deepseek.WithHTTPClient(myHTTPClient),
)
```

## Streaming Chat Completion
Use `CreateChatCompletionStream` when you want tokens as they arrive:

```go
package main

import (
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

	stream, err := client.CreateChatCompletionStream(context.Background(), deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepseekV4Flash,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "Write a short haiku about Go concurrency.",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(resp.Choices) == 0 {
			continue
		}

		fmt.Print(resp.Choices[0].Delta.Content)
	}
}
```

## Function Calling
```go
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

	req := deepseek.ChatCompletionRequest{
		Model: deepseek.DeepseekV4Pro,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "What's the weather in Hangzhou?",
			},
		},
		Tools: []deepseek.Tools{
			{
				Type: "function",
				Function: deepseek.Function{
					Name:        "get_weather",
					Description: "Get weather for a city",
					Parameters: &deepseek.Parameters{
						Type: "object",
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "City name",
							},
						},
						Required: []string{"location"},
					},
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		log.Fatal("no tool call returned")
	}

	call := resp.Choices[0].Message.ToolCalls[0]
	fmt.Printf("tool=%s args=%s\n", call.Function.Name, call.Function.Arguments)
}
```

## Local Models with Ollama
You can also use local models through Ollama:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient(
		"",
		deepseek.WithBaseURL("http://localhost:11434"),
	)

	req := deepseek.OllamaChatRequest{
		Model: "deepseek-r1:7b",
		Messages: []deepseek.OllamaChatMessage{
			{
				Role:    "user",
				Content: "Hello!",
			},
		},
	}

	resp, err := client.CreateOllamaChatCompletion(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Message.Content)
}
```

If you do not have Ollama installed yet:

```bash
ollama run deepseek-r1
```

## Example Index
Runnable examples live in [`examples/`](./examples):

| Scenario | Path |
| --- | --- |
| Chat completion | [`examples/chat`](./examples/chat) |
| Stream chat completion | [`examples/chat_stream`](./examples/chat_stream) |
| Function calling | [`examples/function_calling`](./examples/function_calling) |
| FIM completion | [`examples/fin`](./examples/fin) |
| Balance query | [`examples/balance`](./examples/balance) |
| Ollama chat | [`examples/ollama_chat`](./examples/ollama_chat) |
| Ollama generate | [`examples/ollama_generate`](./examples/ollama_generate) |
| Ollama embeddings | [`examples/ollama_embed`](./examples/ollama_embed) |
| Ollama function calling | [`examples/ollama_function_calling`](./examples/ollama_function_calling) |
| Qwen3 stream chat | [`examples/qwen3_chat_stream`](./examples/qwen3_chat_stream) |
| Qwen3 function calling | [`examples/qwen3_function_calling`](./examples/qwen3_function_calling) |
| QwQ stream chat | [`examples/qwq_chat_stream`](./examples/qwq_chat_stream) |
| QwQ function calling | [`examples/qwq_function_calling`](./examples/qwq_function_calling) |

## Notes
- `NewClient(token)` uses DeepSeek's default base URL automatically.
- Plain tokens are sent as `Bearer <token>` automatically.
- For compatible providers, set a custom base URL with `WithBaseURL(...)`.
- For local Ollama usage, no auth token is required.

## Development
Run tests with:

```bash
go test ./...
```

## License
This project is licensed under the [MIT License](./LICENSE).
