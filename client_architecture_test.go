package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClientAppliesDefaultsAndOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		client := NewClient("token")

		if client.BaseUrl != defaultBaseURL {
			t.Fatalf("expected default base url %q, got %q", defaultBaseURL, client.BaseUrl)
		}
		if client.httpClient == nil {
			t.Fatal("expected http client to be initialized")
		}
		if client.httpClient.Timeout != defaultHTTPTimeout {
			t.Fatalf("expected default timeout %s, got %s", defaultHTTPTimeout, client.httpClient.Timeout)
		}
	})

	t.Run("options", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 5 * time.Second}
		client := NewClient(
			"token",
			WithBaseURL(" https://example.com/ "),
			WithHTTPClient(httpClient),
		)

		if client.BaseUrl != "https://example.com" {
			t.Fatalf("expected trimmed base url, got %q", client.BaseUrl)
		}
		if client.httpClient != httpClient {
			t.Fatal("expected custom http client to be reused")
		}
	})
}

func TestCreateChatCompletionUsesSharedInfrastructure(t *testing.T) {
	var authHeader string
	var contentType string
	var requestBody ChatCompletionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		contentType = r.Header.Get("Content-Type")

		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != chatCompletionSuffix {
			t.Fatalf("expected path %q, got %q", chatCompletionSuffix, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"chatcmpl-1","object":"chat.completion","created":1,"model":"deepseek-chat","choices":[{"index":0,"message":{"role":"assistant","content":"hello"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2,"prompt_cache_hit_tokens":0,"prompt_cache_miss_tokens":1},"system_fingerprint":"fp"}`))
	}))
	defer server.Close()

	client := NewClient(
		"test-token",
		WithBaseURL(server.URL+"/"),
		WithHTTPClient(server.Client()),
	)

	resp, err := client.CreateChatCompletion(context.Background(), &ChatCompletionRequest{
		Model: DeepSeekChat,
		Messages: []ChatCompletionMessage{
			{Role: ChatMessageRoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion returned error: %v", err)
	}

	if authHeader != "Bearer test-token" {
		t.Fatalf("expected bearer auth header, got %q", authHeader)
	}
	if contentType != "application/json" {
		t.Fatalf("expected json content type, got %q", contentType)
	}
	if requestBody.Model != DeepSeekChat {
		t.Fatalf("expected model %q, got %q", DeepSeekChat, requestBody.Model)
	}
	if len(resp.Choices) != 1 || resp.Choices[0].Message.Content != "hello" {
		t.Fatalf("unexpected response payload: %+v", resp.Choices)
	}
}

func TestCreateChatCompletionReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"bad request"}}`))
	}))
	defer server.Close()

	client := NewClient(
		"test-token",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)

	_, err := client.CreateChatCompletion(context.Background(), &ChatCompletionRequest{
		Model: DeepSeekChat,
		Messages: []ChatCompletionMessage{
			{Role: ChatMessageRoleUser, Content: "hi"},
		},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Fatalf("expected parsed error message, got %q", apiErr.Message)
	}
}

func TestGetBalanceOmitsJSONBodyOnGET(t *testing.T) {
	var contentType string
	var contentLength int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		contentLength = r.ContentLength

		if r.Method != http.MethodGet {
			t.Fatalf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != BalanceSuffix {
			t.Fatalf("expected path %q, got %q", BalanceSuffix, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"is_available":true,"balance_infos":[{"currency":"CNY","total_balance":"10","granted_balance":"0","topped_up_balance":"10"}]}`))
	}))
	defer server.Close()

	client := NewClient(
		"test-token",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)

	resp, err := client.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("GetBalance returned error: %v", err)
	}

	if contentType != "" {
		t.Fatalf("expected empty content type for GET without body, got %q", contentType)
	}
	if contentLength != 0 {
		t.Fatalf("expected empty GET body, got content length %d", contentLength)
	}
	if !resp.IsAvailable || len(resp.BalanceInfos) != 1 {
		t.Fatalf("unexpected balance response: %+v", resp)
	}
}
