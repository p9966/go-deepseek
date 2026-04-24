package deepseek

import (
	"context"
	"errors"
	"net/http"
)

var ollamaGenerateSuffix = "/api/generate"

type OllamaGenerateRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`               // 可选：用于生成响应的提示词
	Stream    bool     `json:"stream"`               // 可选：为 false 时返回单个完整响应，而不是流式响应
	Suffix    string   `json:"suffix,omitempty"`     // 可选：模型响应后的后缀文本
	Images    []string `json:"images,omitempty"`     // 可选：Base64 编码图片列表，适用于 llava 等多模态模型
	Format    any      `json:"format,omitempty"`     // 可选：返回格式，可为 json 或 JSON Schema
	Options   *Options `json:"options,omitempty"`    // 可选：模型附加参数，例如 temperature 等 Modelfile 支持的参数
	System    string   `json:"system,omitempty"`     // 可选：系统消息，会覆盖 Modelfile 中定义的 system
	Template  string   `json:"template,omitempty"`   // 可选：提示模板，会覆盖 Modelfile 中定义的 template
	Raw       bool     `json:"raw,omitempty"`        // 可选：为 true 时不对提示词做额外格式化，适合自行传入完整模板
	KeepAlive int      `json:"keep_alive,omitempty"` // 可选：请求结束后模型在内存中的保活时长，默认 5 分钟
}

type Options struct {
	NumKeep          int      `json:"num_keep,omitempty"`
	Seed             int      `json:"seed,omitempty"`
	NumPredict       int      `json:"num_predict,omitempty"`
	TopK             int      `json:"top_k,omitempty"`
	TopP             float64  `json:"top_p,omitempty"`
	MinP             float64  `json:"min_p,omitempty"`
	TypicalP         float64  `json:"typical_p,omitempty"`
	RepeatLastN      int      `json:"repeat_last_n,omitempty"`
	Temperature      float64  `json:"temperature,omitempty"`
	RepeatPenalty    float64  `json:"repeat_penalty,omitempty"`
	PresencePenalty  float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"`
	Mirostat         int      `json:"mirostat,omitempty"`
	MirostatTau      float64  `json:"mirostat_tau,omitempty"`
	MirostatEta      float64  `json:"mirostat_eta,omitempty"`
	PenalizeNewline  bool     `json:"penalize_newline,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Numa             bool     `json:"numa,omitempty"`
	NumCtx           int      `json:"num_ctx,omitempty"`
	NumBatch         int      `json:"num_batch,omitempty"`
	NumGpu           int      `json:"num_gpu,omitempty"`
	MainGpu          int      `json:"main_gpu,omitempty"`
	LowVram          bool     `json:"low_vram,omitempty"`
	VocabOnly        bool     `json:"vocab_only,omitempty"`
	UseMmap          bool     `json:"use_mmap,omitempty"`
	UseMlock         bool     `json:"use_mlock,omitempty"`
	NumThread        int      `json:"num_thread,omitempty"`
}

type OllamaGenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"` // 流式响应时该字段为空，非流式响应时包含完整结果
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`              // 本次响应使用的上下文编码，可在下次请求中继续复用以保持会话记忆
	TotalDuration      int64  `json:"total_duration,omitempty"`       // 生成响应的总耗时
	LoadDuration       int64  `json:"load_duration,omitempty"`        // 加载模型耗时，单位为纳秒
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`    // 提示词 token 数量
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"` // 评估提示词耗时，单位为纳秒
	EvalCount          int    `json:"eval_count,omitempty"`           // 响应 token 数量
	EvalDuration       int64  `json:"eval_duration,omitempty"`        // 生成响应耗时，单位为纳秒
}

func (c *Client) CreateOllamaGenerate(ctx context.Context, req *OllamaGenerateRequest) (*OllamaGenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := c.newRequest(ctx, http.MethodPost, ollamaGenerateSuffix, req)
	if err != nil {
		return nil, err
	}

	var generateResp OllamaGenerateResponse
	if err := c.doJSON(request, &generateResp); err != nil {
		return nil, err
	}

	return &generateResp, nil
}
