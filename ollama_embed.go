package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	deepseek "github.com/p9966/go-deepseek/internal"
)

var ollamaEmbedSuffix = "/api/embed"

type OllamaEmbedRequest struct {
	Model     string   `json:"model"`                // name of model to generate embeddings from
	Input     any      `json:"input"`                // text or list of text to generate embeddings for
	Truncate  bool     `json:"truncate,omitempty"`   // Optional: truncates the end of each input to fit within context length. Returns error if false and context length is exceeded. Defaults to true
	Options   *Options `json:"options,omitempty"`    // Optional: additional model parameters listed in the documentation for the Modelfile such as temperature, https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
	KeepAlive int      `json:"keep_alive,omitempty"` // Optional: controls how long the model will stay loaded into memory following the request (default: 5m)
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

	request, err := deepseek.NewRequestBuilder(c.AuthToken).SetMethod(http.MethodPost).SetBaseUrl(c.BaseUrl).SetPath(ollamaEmbedSuffix).SetBody(req).Build(ctx)
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

	var embedResp OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}
	return &embedResp, nil
}
