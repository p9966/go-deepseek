package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	deepseek "github.com/p9966/go-deepseek/internal"
)

var ollamaChatCompletionSuffix = "/api/chat"

type OllamaChatRequest struct {
	Model     string              `json:"model"`
	Messages  []OllamaChatMessage `json:"messages"`             // the messages of the chat, this can be used to keep a chat memory
	Tools     []Tools             `json:"tools,omitempty"`      // Optional: the tools to use for the chat
	Format    map[string]any      `json:"format,omitempty"`     // Optional: the format to return a response in. Format can be json or a JSON schema
	Stream    bool                `json:"stream"`               // Optional: if false the response will be returned as a single response object, rather than a stream of objects
	Options   *Options            `json:"options,omitempty"`    // Optional: additional model parameters listed in the documentation for the Modelfile such as temperature, https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
	KeepAlive int                 `json:"keep_alive,omitempty"` // Optional: controls how long the model will stay loaded into memory following the request (default: 5m)
}

type OllamaChatMessage struct {
	Role      string       `json:"role"`
	Content   string       `json:"content"`
	Images    []string     `json:"images,omitempty"`
	ToolCalls []OllamaTool `json:"tool_calls,omitempty"`
}

type OllamaTool struct {
	Function struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	} `json:"function"`
}

type OllamaChatResponse struct {
	Model              string             `json:"model"`
	CreatedAt          string             `json:"created_at"`
	Message            *OllamaChatMessage `json:"message"`
	DoneReason         string             `json:"done_reason"`
	Done               bool               `json:"done"`
	TotalDuration      int64              `json:"total_duration"`
	LoadDuration       int64              `json:"load_duration"`
	PromptEvalCount    int                `json:"prompt_eval_count"`
	PromptEvalDuration int64              `json:"prompt_eval_duration"`
	EvalCount          int                `json:"eval_count"`
	EvalDuration       int64              `json:"eval_duration"`
}

func (c *Client) CreateOllamaChatCompletion(ctx context.Context, req *OllamaChatRequest) (*OllamaChatResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := deepseek.NewRequestBuilder(c.AuthToken).SetMethod(http.MethodPost).SetBaseUrl(c.BaseUrl).SetPath(ollamaChatCompletionSuffix).SetBody(req).Build(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}

	var generateResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return nil, err
	}

	return &generateResp, nil
}
