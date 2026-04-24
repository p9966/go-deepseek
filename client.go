package deepseek

import (
	"net/http"
	"strings"
	"time"
)

type Client struct {
	AuthToken  string
	BaseUrl    string
	httpClient *http.Client
}

const (
	defaultBaseURL     = "https://api.deepseek.com"
	defaultHTTPTimeout = 120 * time.Second
)

type ClientOption func(*Client)

// NewClient 使用给定的 API Key 创建一个新的 DeepSeek 客户端。
func NewClient(token string, opts ...ClientOption) *Client {
	client := &Client{
		AuthToken: token,
		BaseUrl:   defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	client.ensureDefaults()
	return client
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseUrl = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout <= 0 {
			return
		}
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.ensureDefaults()
	return c.httpClient.Do(req)
}

func (c *Client) ensureDefaults() {
	if c.BaseUrl == "" {
		c.BaseUrl = defaultBaseURL
	}

	c.BaseUrl = strings.TrimRight(strings.TrimSpace(c.BaseUrl), "/")
	if c.BaseUrl == "" {
		c.BaseUrl = defaultBaseURL
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: defaultHTTPTimeout,
		}
	}
}

func formatAuthToken(token string) string {
	token = strings.TrimSpace(token)
	switch {
	case token == "":
		return ""
	case strings.Contains(token, " "):
		return token
	default:
		return "Bearer " + token
	}
}
