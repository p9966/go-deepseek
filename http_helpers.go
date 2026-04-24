package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	deepseekinternal "github.com/p9966/go-deepseek/internal"
)

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return fmt.Sprintf("deepseek api error: status %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("deepseek api error: status %d", e.StatusCode)
}

type apiErrorEnvelope struct {
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	c.ensureDefaults()

	return deepseekinternal.NewRequestBuilder(formatAuthToken(c.AuthToken)).
		SetMethod(method).
		SetBaseUrl(c.BaseUrl).
		SetPath(path).
		SetBody(body).
		Build(ctx)
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return err
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}

	var envelope apiErrorEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil {
		switch {
		case envelope.Error != nil && envelope.Error.Message != "":
			apiErr.Message = envelope.Error.Message
		case envelope.Message != "":
			apiErr.Message = envelope.Message
		}
	}

	return apiErr
}
