# Go DeepSeek
[![Go Reference](https://pkg.go.dev/badge/github.com/p9966/go-deepseek.svg)](https://pkg.go.dev/github.com/p9966/go-deepseek)
[![Go Report Card](https://goreportcard.com/badge/github.com/p9966/go-deepseek)](https://goreportcard.com/report/github.com/p9966/go-deepseek)


This library provides unofficial Go clients for [DeepSeek](https://www.deepseek.com/). We support: 

* deepseek-chat

## Installation

```
go get github.com/p9966/go-deepseek
```
Currently, go-deepseek requires Go version 1.23.5 or greater.


## Usage

### example usage:

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

### Getting an DeepSeek API Key:

1. Visit the DeepSeek website at [https://www.deepseek.com/](https://www.deepseek.com/).
2. If you don't have an account, click on "Sign Up" to create one. If you do, click "Log In".
3. Once logged in, navigate to your API key management page.
4. Click on "Create new secret key".
5. Enter a name for your new key, then click "Create secret key".
6. Your new API key will be displayed. Use this key to interact with the DeepSeek API.

**Note:** Your API key is sensitive information. Do not share it with anyone.