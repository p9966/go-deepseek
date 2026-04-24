package deepseek

import (
	"context"
	"errors"
	"net/http"
)

var ollamaChatCompletionSuffix = "/api/chat"

type OllamaChatRequest struct {
	Model     string              `json:"model"`
	Messages  []OllamaChatMessage `json:"messages"`             // 聊天消息列表，可用于维持会话上下文
	Tools     []Tools             `json:"tools,omitempty"`      // 可选：本次聊天可用的工具列表
	Format    map[string]any      `json:"format,omitempty"`     // 可选：返回格式，可为 json 或 JSON Schema
	Stream    bool                `json:"stream"`               // 可选：为 false 时返回单个完整响应，而不是流式响应
	Options   *Options            `json:"options,omitempty"`    // 可选：模型附加参数，例如 temperature 等 Modelfile 支持的参数
	KeepAlive int                 `json:"keep_alive,omitempty"` // 可选：请求结束后模型在内存中的保活时长，默认 5 分钟
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

	request, err := c.newRequest(ctx, http.MethodPost, ollamaChatCompletionSuffix, req)
	if err != nil {
		return nil, err
	}

	var generateResp OllamaChatResponse
	if err := c.doJSON(request, &generateResp); err != nil {
		return nil, err
	}

	return &generateResp, nil
}
