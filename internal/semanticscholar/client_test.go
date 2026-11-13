package semanticscholar

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 429,
		Message:    "Too Many Requests",
	}

	expected := "semantic scholar api error: status=429, message=Too Many Requests"

	if err.Error() != expected {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			err.Error(),
		)
	}
}

func TestGet(t *testing.T) {
	// Create a fake API server.
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// Check HTTP method.
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			// Check API key.
			if got := r.Header.Get("x-api-key"); got != "test-api-key" {
				t.Errorf("unexpected API key: %s", got)
			}

			// Check query parameter.
			if got := r.URL.Query().Get("query"); got != "golang" {
				t.Errorf("unexpected query: %s", got)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "Hello from fake API",
			})
		},
	))
	defer server.Close()

	// Create client.
	client := NewClinet("test-api-key")
	client.baseURL = server.URL

	// Prepare request.
	params := url.Values{}
	params.Set("query", "golang")

	var result struct {
		Message string `json:"message"`
	}

	// Execute request.
	err := client.get(
		context.Background(),
		"/paper/search",
		params,
		&result,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response.
	if result.Message != "Hello from fake API" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestGetAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "rate limit",
			statusCode: http.StatusTooManyRequests,
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					http.Error(
						w,
						"API error",
						tt.statusCode,
					)
				},
			))
			defer server.Close()

			client := NewClinet("")
			client.baseURL = server.URL

			var result map[string]any

			err := client.get(
				context.Background(),
				"/paper/search",
				nil,
				&result,
			)

			if err == nil {
				t.Fatal("expected an API error, got nil")
			}

			var apiErr *APIError

			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *APIError, got %T", err)
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf(
					"expected status %d, got %d",
					tt.statusCode,
					apiErr.StatusCode,
				)
			}
		})
	}
}

func TestGetContextCancellation(t *testing.T) {
	started := make(chan struct{})

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(started)

			<-r.Context().Done()
		}),
	)

	defer server.Close()

	client := NewClinet("")
	client.baseURL = server.URL

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	errCh := make(chan error, 1)

	go func() {

		var result map[string]any

		errCh <- client.get(
			ctx,
			"/paper/search",
			nil,
			&result,
		)
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatalf("server did not recive the rquest")
	}

	cancel()

	err := <-errCh

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Cancel, got: %v", err)
	}
}
