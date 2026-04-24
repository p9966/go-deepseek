package deepseek

import (
	"context"
	"errors"
	"net/http"
)

var ollamaEmbedSuffix = "/api/embed"

type OllamaEmbedRequest struct {
	Model     string   `json:"model"`                // 用于生成向量的模型名称
	Input     any      `json:"input"`                // 需要生成向量的文本，支持单条或多条
	Truncate  bool     `json:"truncate,omitempty"`   // 可选：是否在超出上下文长度时截断输入，默认开启；关闭时若超限将返回错误
	Options   *Options `json:"options,omitempty"`    // 可选：模型附加参数，例如 temperature 等 Modelfile 支持的参数
	KeepAlive int      `json:"keep_alive,omitempty"` // 可选：请求结束后模型在内存中的保活时长，默认 5 分钟
}

type OllamaEmbedResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float64 `json:"embeddings"`
	TotalDuration   int64       `json:"total_duration"`
	LoadDuration    int64       `json:"load_duration"`
	PromptEvalCount int         `json:"prompt_eval_count"`
}

func (c *Client) CreateOllamaEmbed(ctx context.Context, req *OllamaEmbedRequest) (*OllamaEmbedResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := c.newRequest(ctx, http.MethodPost, ollamaEmbedSuffix, req)
	if err != nil {
		return nil, err
	}

	var embedResp OllamaEmbedResponse
	if err := c.doJSON(request, &embedResp); err != nil {
		return nil, err
	}
	return &embedResp, nil
}
