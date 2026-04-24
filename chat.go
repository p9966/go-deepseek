// Package deepseek 提供 DeepSeek API 的非官方 Go 客户端。
// 当前支持与 DeepSeek 的聊天能力进行交互。
package deepseek

import (
	"context"
	"errors"
	"net/http"
)

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

const chatCompletionSuffix = "/chat/completions"

type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	FrequencyPenalty float32                 `json:"frequency_penalty"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`       // 可选：最大输出 token 数，需大于 1
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"` // 可选：出现惩罚，取值范围为 [-2, 2]
	Temperature      float32                 `json:"temperature,omitempty"`      // 可选：采样温度，最大为 2
	TopP             float32                 `json:"top_p,omitempty"`            // 可选：核采样参数，最大为 1
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`  // 可选：自定义响应格式
	Stop             []string                `json:"stop,omitempty"`             // 可选：停止生成的标记
	Tools            []Tools                 `json:"tools,omitempty"`            // 可选：工具列表
	LogProbs         bool                    `json:"logprobs,omitempty"`         // 可选：是否返回对数概率
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`     // 可选：返回对数概率最高的 token 数量，最大为 20
	ToolChoice       string                  `json:"tool_choice,omitempty"`      // 可选：工具选择策略，可为 "auto" 或 "none"
	ExtraBody        ExtraBody               `json:"extra_body,omitempty"`       // 可选：额外请求体，用于思考相关配置
	ReasoningEffort  ReasoningEffortModel    `json:"reasoning_effort,omitempty"` // 可选：思考强度，可为 "high" 或 "max"
}

type ExtraBody struct {
	Thinking ExtraBodyThinking `json:"thinking,omitempty"`
}

type ExtraBodyThinking struct {
	Type ExtraBodyThinkingModel `json:"type,omitempty"`
}

type ReasoningEffortModel string

const (
	ReasoningEffortModelHigh ReasoningEffortModel = "high"
	ReasoningEffortModelMax  ReasoningEffortModel = "max"
)

type ExtraBodyThinkingModel string

const (
	ExtraBodyThinkingModelEnabled  ExtraBodyThinkingModel = "enabled"
	ExtraBodyThinkingModelDisabled ExtraBodyThinkingModel = "disabled"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Tools struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string      `json:"name"`                 // 必填：函数名称
	Description string      `json:"description"`          // 必填：函数描述
	Parameters  *Parameters `json:"parameters,omitempty"` // 可选：参数结构定义
}

type Parameters struct {
	Type       string                 `json:"type"` // 必填：参数类型，例如 "object"
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

type ChatCompletionResponse struct {
	ID                string   `json:"id"`                 // 聊天补全结果的唯一标识
	Object            string   `json:"object"`             // 对象类型，通常为 "chat.completion"
	Created           int64    `json:"created"`            // 聊天补全创建时间戳
	Model             string   `json:"model"`              // 生成该结果所使用的模型
	Choices           []Choice `json:"choices"`            // 模型生成的候选结果列表
	Usage             Usage    `json:"usage"`              // token 使用统计信息
	SystemFingerprint string   `json:"system_fingerprint"` // 系统配置指纹
}

type Choice struct {
	Index        int       `json:"index"`              // 当前候选结果在列表中的索引
	Message      Message   `json:"message"`            // 模型生成的消息内容
	LogProbs     *LogProbs `json:"logprobs,omitempty"` // token 的对数概率信息
	FinishReason string    `json:"finish_reason"`      // 生成结束原因
}

type Message struct {
	Role             string     `json:"role"`                        // 消息发送者角色，例如 "user"、"assistant"
	Content          string     `json:"content"`                     // 消息内容
	ReasoningContent string     `json:"reasoning_content,omitempty"` // 可选：思考过程内容
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // 可选：工具调用列表
}

type ToolCall struct {
	Index    int          `json:"index"` // 当前工具调用在列表中的索引
	Id       string       `json:"id"`    // 工具调用 ID
	Type     string       `json:"type"`  // 工具调用类型，例如 "function"
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`      // 要调用的函数名称
	Arguments string `json:"arguments"` // 函数调用参数内容
}

type LogProbs struct {
	Tokens        []string             `json:"tokens,omitempty"`         // token 列表
	TokenLogProbs []float64            `json:"token_logprobs,omitempty"` // 每个 token 的对数概率
	TopLogProbs   []map[string]float64 `json:"top_logprobs,omitempty"`   // 每个 token 对应的最高对数概率集合
}

type Usage struct {
	PromptTokens          int `json:"prompt_tokens"`            // 提示词使用的 token 数
	CompletionTokens      int `json:"completion_tokens"`        // 补全内容使用的 token 数
	TotalTokens           int `json:"total_tokens"`             // 总 token 数
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`  // 命中缓存的 token 数
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"` // 未命中缓存的 token 数
}

func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := c.newRequest(ctx, http.MethodPost, chatCompletionSuffix, req)
	if err != nil {
		return nil, err
	}

	var result ChatCompletionResponse
	if err := c.doJSON(request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
