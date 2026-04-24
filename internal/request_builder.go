package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type AuthRequest struct {
	baseUrl   string
	authToken string
	path      string
	method    string
	body      any
}

type RequestBuilder interface {
	SetBaseUrl(string) *AuthRequest
	SetBody(any) *AuthRequest
	SetPath(string) *AuthRequest
	SetMethod(string) *AuthRequest
	Build(context.Context) (*http.Request, error)
}

func NewRequestBuilder(authToken string) *AuthRequest {
	return &AuthRequest{
		authToken: authToken,
	}
}

func (r *AuthRequest) SetBaseUrl(baseUrl string) *AuthRequest {
	r.baseUrl = baseUrl
	return r
}

func (r *AuthRequest) SetPath(path string) *AuthRequest {
	r.path = path
	return r
}

func (r *AuthRequest) SetMethod(method string) *AuthRequest {
	r.method = method
	return r
}

func (r *AuthRequest) SetBody(body any) *AuthRequest {
	r.body = body
	return r
}

func (r *AuthRequest) Build(ctx context.Context) (*http.Request, error) {
	var bodyReader io.Reader
	switch v := r.body.(type) {
	case nil:
		bodyReader = nil
	case io.Reader:
		bodyReader = v
	default:
		var reqBody []byte
		reqBody, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqBody)
	}

	if r.baseUrl == "" {
		return nil, errors.New("base url can not be empty")
	}
	if r.method == "" {
		return nil, errors.New("method can not be empty")
	}

	endpoint := strings.TrimRight(r.baseUrl, "/")
	if r.path != "" {
		if !strings.HasPrefix(r.path, "/") {
			endpoint += "/"
		}
		endpoint += r.path
	}

	req, err := http.NewRequestWithContext(ctx, r.method, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	if r.authToken != "" {
		req.Header.Set("Authorization", r.authToken)
	}
	if r.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
