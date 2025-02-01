# Go DeepSeek
[![Go Reference](https://pkg.go.dev/badge/github.com/p9966/go-deepseek.svg)](https://pkg.go.dev/github.com/p9966/go-deepseek)
[![Go Report Card](https://goreportcard.com/badge/github.com/p9966/go-deepseek)](https://goreportcard.com/report/github.com/p9966/go-deepseek)

This library provides unofficial Go clients for [DeepSeek](https://www.deepseek.com/). It currently supports: 

* Chat Completion
* FIM Completion
* Function Calling

## Installation
To install the library, run:

```
go get github.com/p9966/go-deepseek
```
**Note:** This library requires Go version 1.23 or higher.


## Usage
### Quick Start:
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

## Obtaining a DeepSeek API Key

1. Visit the DeepSeek website [DeepSeek website](https://www.deepseek.com/).
2. Sign up or log in to your account.
3. Navigate to the API key management page.
4. Click **"Create new secret key"**.
5. Enter a name for your key and confirm.
6. Your API key will be displayed. Use it to interact with the DeepSeek API.

**Important:** Keep your API key secure and do not share it publicly.

## Other examples:
<details>
<summary>FIM Completion</summary>

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
```
</details>

## License
This project is licensed under the [MIT License](LICENSE).