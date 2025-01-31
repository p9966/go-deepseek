package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type AuthRequest struct {
	baseUrl   string
	authToken string
	path      string
	body      any
}

type RequestBuilder interface {
	SetBaseUrl(string) *AuthRequest
	SetBody([]byte) *AuthRequest
	SetPath(string) *AuthRequest
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

func (r *AuthRequest) SetBody(body any) *AuthRequest {
	r.body = body
	return r
}

func (r *AuthRequest) Build(ctx context.Context) (*http.Request, error) {
	var bodyReader io.Reader
	if v, ok := r.body.(io.Reader); ok {
		bodyReader = v
	} else {
		var reqBody []byte
		reqBody, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqBody)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseUrl+r.path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.authToken)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
