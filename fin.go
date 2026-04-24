package deepseek

import (
	"context"
	"errors"
	"net/http"
)

const finCompletionSuffix = "/beta/completions"

type FINCompletionRequest struct {
	Model            string    `json:"model"`                       // 必填：模型 ID
	Prompt           string    `json:"prompt"`                      // 必填：用于生成补全内容的提示词
	Echo             bool      `json:"echo,omitempty"`              // 可选：是否返回输入提示词
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"` // 可选：控制生成内容重复度，取值范围 [-2, 2]
	Logprobs         int       `json:"logprobs,omitempty"`          // 可选：返回最可能输出 token 的对数概率，同时包含已采样 token 的对数概率，最大值为 20
	MaxTokens        int       `json:"max_tokens,omitempty"`        // 可选：生成内容的最大长度
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`  // 可选：控制生成内容多样性，取值范围 [-2, 2]
	Stop             *[]string `json:"stop,omitempty"`              // 可选：停止生成的字符串或字符串数组
	Stream           bool      `json:"stream,omitempty"`            // 可选：是否以流式方式返回结果
	Suffix           *string   `json:"suffix,omitempty"`            // 可选：被补全内容的后缀
	Temperature      float64   `json:"temperature,omitempty"`       // 可选：采样温度，建议与 top_p 二选一调整，不建议同时修改
	TopP             float64   `json:"top_p,omitempty"`             // 可选：核采样参数，建议与 temperature 二选一调整，不建议同时修改
}

type FINCompletionResponse struct {
	ID                string                `json:"id"`
	Choices           []FINCompletionChoice `json:"choices"`
	Created           int                   `json:"created"`
	Model             string                `json:"model"`
	SystemFingerprint string                `json:"system_fingerprint"`
	Object            string                `json:"object"`
	Usage             struct {
		CompletionTokens      int                                   `json:"completion_tokens"`
		PromptTokens          int                                   `json:"prompt_tokens"`
		PromptCacheHitTokens  int                                   `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int                                   `json:"prompt_cache_miss_tokens"`
		TotalTokens           int                                   `json:"total_tokens"`
		PromptTokensDetails   FINCompletionUsagePromptTokensDetails `json:"prompt_tokens_details"`
	} `json:"usage"`
}

type FINCompletionChoice struct {
	FinishReason string                      `json:"finish_reason"`
	Index        int                         `json:"index"`
	Logprobs     FINCompletionChoiceLogprobs `json:"logprobs"`
	Text         string                      `json:"text"`
}

type FINCompletionUsagePromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type FINCompletionChoiceLogprobs struct {
	Tokens        []string  `json:"tokens"`
	TokenLogprobs []float64 `json:"token_logprobs"`
}

func (c *Client) CreateFINCompletion(ctx context.Context, req *FINCompletionRequest) (*FINCompletionResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := c.newRequest(ctx, http.MethodPost, finCompletionSuffix, req)
	if err != nil {
		return nil, err
	}

	var finResponse FINCompletionResponse
	if err := c.doJSON(request, &finResponse); err != nil {
		return nil, err
	}

	return &finResponse, nil
}
