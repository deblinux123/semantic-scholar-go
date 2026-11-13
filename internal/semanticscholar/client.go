package semanticscholar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL        = "https://api.semanticscholar.org/graph/v1"
	DefaultTimeout = 40 * time.Second
)

type Client struct {
	httpClinet *http.Client
	apiKey     string
	baseURL    string
}

func NewClinet(apiKey string) *Client {
	return &Client{
		httpClinet: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiKey:  apiKey,
		baseURL: BaseURL,
	}
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf(
		"semantic scholar api error: status=%d, message=%s",
		e.StatusCode,
		e.Message,
	)
}

func (c *Client) get(
	ctx context.Context,
	endpoint string,
	params url.Values,
	result any,
) error {
	requestURL := c.baseURL + endpoint

	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
		nil,
	)

	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(
		"Accept",
		"application/json",
	)

	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClinet.Do(req)

	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(
			io.LimitReader(resp.Body, 4096),
		)

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
