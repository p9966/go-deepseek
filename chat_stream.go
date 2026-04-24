package deepseek

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StreamChatCompletionRequest struct {
	Stream           bool                    `json:"stream"`
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	FrequencyPenalty float32                 `json:"frequency_penalty"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	Tools            []Tools                 `json:"tools,omitempty"`
	LogProbs         bool                    `json:"logprobs,omitempty"`
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`
	EnableThink      bool                    `json:"enable_thinking"`
}

type StreamChatCompletionResponse struct {
	ID                string              `json:"id"`
	Object            string              `json:"object"`
	Created           int64               `json:"created"`
	Model             string              `json:"model"`
	Choices           []StreamChatChoices `json:"choices"`
	SystemFingerprint string              `json:"system_fingerprint"`
}

type StreamChatChoices struct {
	Index        int                  `json:"index"`
	Delta        StreamChatChoiceData `json:"delta"`
	LogProbs     *LogProbs            `json:"logprobs,omitempty"`
	FinishReason string               `json:"finish_reason"`
}

type StreamChatChoiceData struct {
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content"`
	ToolCalls        []ToolCall `json:"tool_calls"`
}

type ChatCompletionStream interface {
	Recv() (*StreamChatCompletionResponse, error)
	Close() error
}

type chatCompletionStream struct {
	ctx    context.Context
	cancel context.CancelFunc
	resp   *http.Response
	reader *bufio.Reader
}

func (s *chatCompletionStream) Recv() (*StreamChatCompletionResponse, error) {
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("failed to read stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "data: [DONE]" {
			return nil, io.EOF
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var resp StreamChatCompletionResponse
		if err := json.Unmarshal([]byte(line[6:]), &resp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return &resp, nil
	}
}

func (s *chatCompletionStream) Close() error {
	s.cancel()
	return s.resp.Body.Close()
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, req StreamChatCompletionRequest) (ChatCompletionStream, error) {
	req.Stream = true
	request, err := c.newRequest(ctx, http.MethodPost, chatCompletionSuffix, req)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &chatCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}, nil
}
