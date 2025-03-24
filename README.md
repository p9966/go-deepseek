# Go DeepSeek
[![Go Reference](https://pkg.go.dev/badge/github.com/p9966/go-deepseek.svg)](https://pkg.go.dev/github.com/p9966/go-deepseek)
[![Go Report Card](https://goreportcard.com/badge/github.com/p9966/go-deepseek)](https://goreportcard.com/report/github.com/p9966/go-deepseek)

This library provides an unofficial Go client for [DeepSeek](https://www.deepseek.com/),it also supports [QwQ](https://help.aliyun.com/zh/model-studio/getting-started/what-is-model-studio), [OpenAI](https://platform.openai.com/docs/overview).enabling interaction with both online and local models. It supports the following features: 
* Chat Completion
* Stream Chat Completion
* FIM (Fill-in-Middle) Completion
* Function Calling
* API balance query
* Embeddings

## Installation
To install the library, run:

```
go get github.com/p9966/go-deepseek
```
**Note:** This library requires Go version 1.23 or higher.


## Usage
### Quick Start:
#### Chat Completion with DeepSeek API
Here’s an example of how to use the library for chat completion:
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
	// client := deepseek.Client{
	//	BaseUrl:   "https://api.deepseek.ai/v1", // deepseek
    //	BaseUrl:   "https://dashscope.aliyuncs.com/compatible-mode/v1",  // qwq
	//	BaseUrl:   "https://api.openai.com/v1/",// openai
	// 	AuthToken: os.Getenv("DEEPSEEK_API_KEY"),
	// }
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

	if len(resp.Choices) == 0 {
		log.Fatal("No response choices available")
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

```

#### Stream Chat Completion with DeepSeek API
Here’s an example of how to use the library for stream chat completion:
```go
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	// client := deepseek.Client{
	// 	BaseUrl:   "https://api.deepseek.ai/v1", // Custom API addresses
	// 	AuthToken: os.Getenv("DEEPSEEK_API_KEY"),
	// }
	scanner := bufio.NewScanner(os.Stdin)

	var messages []deepseek.ChatCompletionMessage

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
			Model:    deepseek.DeepSeekChat,
			Messages: messages,
		}

		ctx := context.Background()
		stream, err := client.CreateChatCompletionStream(ctx, request)
		if err != nil {
			log.Fatalf("ChatCompletionStream failed: %v", err)
		}
		defer stream.Close()

		fmt.Print("DeepSeek: ")
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
```

#### Local Model via Ollama
To use a local model with Ollama:
```go
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
	request := deepseek.OllamaChatRequest{
		Model: "deepseek-r1:7b",
		Messages: []deepseek.OllamaChatMessage{
			{
				Role:    "user",
				Content: "Hello!",
			},
		},
	}
	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(response.Message.Content)
}

```

## Obtaining a DeepSeek API Key

1. Visit the DeepSeek website [DeepSeek website](https://www.deepseek.com/).
2. Sign up or log in to your account.
3. Navigate to the API key management page.
4. Click **"Create new secret key"**.
5. Enter a name for your key and confirm.
6. Your API key will be displayed. Use it to interact with the DeepSeek API.

**Important:** Keep your API key secure and do not share it publicly.

## Local Model Installation
1. Visit the [Ollama website](https://ollama.com/).
2. Download and install Ollama.
3. Open a terminal and run the following command to download the model:
	```bash
	ollama run deepseek-r1
	```
4. You can now use the model locally.

**Note:** You can also download other models from the[Ollama model library](https://ollama.com/search) and use them in the same way.

## Other examples:
<details>
<summary>FIM (Fill-in-Middle) Completion</summary>

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

```
</details>

<details>
<summary>Function Calling</summary>

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
	request := deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: "How's the weather in Hangzhou?",
			},
		},
		Tools: []deepseek.Tools{
			{
				Type: "function",
				Function: deepseek.Function{
					Name:        "get_weather",
					Description: "Get weather of an location, the user shoud supply a location first",
					Parameters: &deepseek.Parameters{
						Type: "object",
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"description": "The location to get weather",
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
	resp, err := client.CreateChatCompletion(ctx, &request)
	if err != nil {
		log.Fatalf("ChatCompletion failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("No response choices available")
	}

	if len(resp.Choices[0].Message.ToolCalls) == 0 {
		log.Fatal("No function calls available")
	}

	fmt.Printf("Function name: %v, args:%s\n", resp.Choices[0].Message.ToolCalls[0].Function.Name, resp.Choices[0].Message.ToolCalls[0].Function.Arguments)
}



// ---------------------qwq---------------------
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

```
</details>

<details>
<summary>Embeddings</summary>

```go
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

```
</details>

<details>
<summary>API balance query</summary>

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
	if resp, err := client.GetBalance(context.TODO()); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.IsAvailable)
	}
}


```
</details>

## Troubleshooting
### Unable to Access DeepSeek API Platform
If you encounter a 503 error when trying to access DeepSeek’s API platform:
1. Check if DeepSeek’s services are down or under maintenance.
2. Clear your browser cache and cookies.
3. Try accessing the platform from a different network or device.
4. Contact DeepSeek’s support team for further assistance. service@deepseek.com

## License
This project is licensed under the [MIT License](LICENSE).