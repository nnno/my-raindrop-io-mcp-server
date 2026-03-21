package raindrop

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

const (
	defaultBaseURL   = "https://api.raindrop.io/rest/v1"
	maxResponseBytes = 10 * 1024 * 1024 // 10MB
)

type ClientOption func(*Client)

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: defaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read up to maxResponseBytes + 1 to detect truncation
	limited := io.LimitReader(resp.Body, maxResponseBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if int64(len(data)) > maxResponseBytes {
		return nil, fmt.Errorf("response too large (exceeds %d bytes)", maxResponseBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyHTTPError(resp.StatusCode, data)
	}

	return data, nil
}

func classifyHTTPError(status int, body []byte) *entity.DomainError {
	msg := fmt.Sprintf("API error (status %d): %s", status, string(body))
	switch {
	case status == http.StatusNotFound:
		return &entity.DomainError{Kind: entity.ErrNotFound, Message: msg}
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return &entity.DomainError{Kind: entity.ErrUnauthorized, Message: msg}
	case status == http.StatusTooManyRequests:
		return &entity.DomainError{Kind: entity.ErrRateLimited, Message: msg}
	default:
		return &entity.DomainError{Kind: entity.ErrInternal, Message: msg}
	}
}
