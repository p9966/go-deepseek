package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	deepseek "github.com/p9966/go-deepseek/internal"
)

const finCompletionSuffix = "/beta/completions"

type FINCompletionRequest struct {
	Model            string    `json:"model"`                       // Required, 模型的 ID
	Prompt           string    `json:"prompt"`                      // Required, 用于生成完成内容的提示
	Echo             bool      `json:"echo,omitempty"`              // Optional, 是否返回输入的提示内容
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"` // Optional, 控制生成内容的重复性，取值范围 [-2, 2]
	Logprobs         int       `json:"logprobs,omitempty"`          // Optional, 制定输出中包含 logprobs 最可能输出 token 的对数概率，包含采样的 token。例如，如果 logprobs 是 20，API 将返回一个包含 20 个最可能的 token 的列表。API 将始终返回采样 token 的对数概率，因此响应中可能会有最多 logprobs+1 个元素。logprobs 的最大值是 20
	MaxTokens        int       `json:"max_tokens,omitempty"`        // Optional, 生成内容的最大长度
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`  // Optional, 控制生成内容的多样性，取值范围 [-2, 2]
	Stop             *[]string `json:"stop,omitempty"`              // Optional, 停止生成内容的字符串或字符串数组
	Stream           bool      `json:"stream,omitempty"`            // Optional, 是否流式返回结果
	Suffix           *string   `json:"suffix,omitempty"`            // Optional, 制定被补全内容的后缀。
	Temperature      float64   `json:"temperature,omitempty"`       // Optional, 采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
	TopP             float64   `json:"top_p,omitempty"`             // Optional, 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
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

	request, err := deepseek.NewRequestBuilder(c.AuthToken).SetBaseUrl(c.BaseUrl).SetPath(finCompletionSuffix).SetBody(req).Build(ctx)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var finResponse FINCompletionResponse
	if err := json.Unmarshal(buf, &finResponse); err != nil {
		return nil, err
	}
	return &finResponse, nil
}
