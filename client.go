package deepseek

import (
	"net/http"
	"time"
)

type Client struct {
	AuthToken  string
	BaseUrl    string
	httpClient *http.Client
}

// NewClient creates a new DeepSeek client with the provided API key.
func NewClient(token string) *Client {
	return &Client{
		AuthToken: token,
		BaseUrl:   "https://api.deepseek.com",
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.httpClient = &http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
