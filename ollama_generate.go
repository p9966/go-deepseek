package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	deepseek "github.com/p9966/go-deepseek/internal"
)

var ollamaGenerateSuffix = "/api/generate"

type OllamaGenerateRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`               // Optional: the prompt to generate a response for
	Stream    bool     `json:"stream"`               // Optional: if false the response will be returned as a single response object, rather than a stream of objects
	Suffix    string   `json:"suffix,omitempty"`     // Optional: the text after the model response
	Images    []string `json:"images,omitempty"`     // Optional: a list of base64-encoded images (for multimodal models such as llava)
	Format    any      `json:"format,omitempty"`     // Optional: the format to return a response in. Format can be json or a JSON schema
	Options   *Options `json:"options,omitempty"`    // Optional: additional model parameters listed in the documentation for the Modelfile such as temperature, https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
	System    string   `json:"system,omitempty"`     // Optional: system message to (overrides what is defined in the Modelfile, https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values)
	Template  string   `json:"template,omitempty"`   // Optional: the prompt template to use (overrides what is defined in the Modelfile, https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values)
	Raw       bool     `json:"raw,omitempty"`        // Optional: if true no formatting will be applied to the prompt. You may choose to use the raw parameter if you are specifying a full templated prompt in your request to the API
	KeepAlive int      `json:"keep_alive,omitempty"` // Optional: controls how long the model will stay loaded into memory following the request (default: 5m)
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
	Response           string `json:"response"` // empty if the response was streamed, if not streamed, this will contain the full response
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`              // an encoding of the conversation used in this response, this can be sent in the next request to keep a conversational memory
	TotalDuration      int64  `json:"total_duration,omitempty"`       // time spent generating the response
	LoadDuration       int64  `json:"load_duration,omitempty"`        // time spent in nanoseconds loading the model
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`    // number of tokens in the prompt
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"` // time spent in nanoseconds evaluating the prompt
	EvalCount          int    `json:"eval_count,omitempty"`           // number of tokens in the response
	EvalDuration       int64  `json:"eval_duration,omitempty"`        // time in nanoseconds spent generating the response
}

func (c *Client) CreateOllamaGenerate(ctx context.Context, req *OllamaGenerateRequest) (*OllamaGenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	request, err := deepseek.NewRequestBuilder(c.AuthToken).SetMethod(http.MethodPost).SetBaseUrl(c.BaseUrl).SetPath(ollamaGenerateSuffix).SetBody(req).Build(ctx)
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

	var generateResp OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return nil, err
	}

	return &generateResp, nil
}
